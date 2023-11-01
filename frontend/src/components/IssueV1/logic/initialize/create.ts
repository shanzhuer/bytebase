import { groupBy, orderBy } from "lodash-es";
import { v4 as uuidv4 } from "uuid";
import { reactive } from "vue";
import { rolloutServiceClient } from "@/grpcweb";
import { TemplateType } from "@/plugins";
import {
  useChangelistStore,
  useCurrentUserV1,
  useDatabaseV1Store,
  useEnvironmentV1Store,
  useProjectV1Store,
  useSheetV1Store,
} from "@/store";
import { getProjectAndSchemaDesignSheetId } from "@/store/modules/v1/common";
import {
  ComposedProject,
  emptyIssue,
  TaskTypeListWithStatement,
  UNKNOWN_ID,
} from "@/types";
import { IssueStatus, Issue_Type } from "@/types/proto/v1/issue_service";
import {
  Plan,
  Plan_ChangeDatabaseConfig,
  Plan_ChangeDatabaseConfig_Type,
  Plan_Spec,
  Plan_Step,
  Rollout,
  Stage,
  Task_Type,
} from "@/types/proto/v1/rollout_service";
import {
  extractSheetUID,
  generateSQLForChangeToDatabase,
  getSheetStatement,
  setSheetNameForTask,
  setSheetStatement,
  sheetNameOfTaskV1,
} from "@/utils";
import { nextUID } from "../base";
import { sheetNameForSpec } from "../plan";
import { createEmptyLocalSheet, getLocalSheetByName } from "../sheet";
import { trySetDefaultAssignee } from "./assignee";

export type InitialSQL = {
  sqlList?: string[];
  sql?: string;
};

type CreateIssueParams = {
  databaseUIDList: string[];
  project: ComposedProject;
  query: Record<string, string>;
  initialSQL: InitialSQL;
  branch?: string;
};

export const createIssueSkeleton = async (query: Record<string, string>) => {
  const issue = emptyIssue();
  const me = useCurrentUserV1();
  issue.creator = `users/${me.value.email}`;
  issue.creatorEntity = me.value;

  const project = await useProjectV1Store().getOrFetchProjectByUID(
    query.project
  );
  issue.project = project.name;
  issue.projectEntity = project;
  issue.uid = nextUID();
  issue.name = `${project.name}/issues/${issue.uid}`;
  issue.title = query.name;
  issue.type = Issue_Type.DATABASE_CHANGE;
  issue.status = IssueStatus.OPEN;

  const databaseUIDList = (query.databaseList ?? "")
    .split(",")
    .filter((uid) => uid && uid !== String(UNKNOWN_ID));
  await prepareDatabaseList(databaseUIDList, project.uid);
  const branch = query.branch || undefined;

  const params: CreateIssueParams = {
    databaseUIDList,
    project,
    query,
    initialSQL: extractInitialSQLListFromQuery(query),
    branch,
  };

  // Prepare params context for building plan.
  await prepareParamsContext(params);

  const plan = await buildPlan(params);
  issue.plan = plan.name;
  issue.planEntity = plan;

  const rollout = await previewPlan(plan, params);
  issue.rollout = rollout.name;
  issue.rolloutEntity = rollout;

  await trySetDefaultAssignee(issue);

  const description = query.description;
  if (description) {
    issue.description = description;
  }

  return issue;
};

const prepareParamsContext = async (params: CreateIssueParams) => {
  if (params.branch) {
    const [_, sheetId] = getProjectAndSchemaDesignSheetId(params.branch);
    await useSheetV1Store().getOrFetchSheetByUID(sheetId);
  }
};

export const buildPlan = async (params: CreateIssueParams) => {
  const { databaseUIDList, project, query } = params;

  const plan = Plan.fromJSON({
    uid: nextUID(),
  });
  plan.name = `${project.name}/plans/${plan.uid}`;
  if (query.changelist) {
    // build plan for changelist
    plan.steps = await buildStepsViaChangelist(
      databaseUIDList,
      query.changelist,
      params
    );
  } else if (query.mode === "tenant") {
    // in tenant mode, all specs share a unique sheet
    const sheetUID = nextUID();
    // build tenant plan
    if (databaseUIDList.length === 0) {
      // evaluate DeploymentConfig and generate steps/specs
      if (query.databaseGroupName) {
        plan.steps = await buildStepsForDatabaseGroup(
          params,
          query.databaseGroupName
        );
      } else {
        plan.steps = await buildStepsViaDeploymentConfig(params, sheetUID);
      }
    } else {
      plan.steps = await buildSteps(databaseUIDList, params, sheetUID);
    }
  } else {
    // build standard plan
    plan.steps = await buildSteps(databaseUIDList, params, undefined);
  }
  return plan;
};

export const buildSteps = async (
  databaseUIDList: string[],
  params: CreateIssueParams,
  sheetUID?: string // if specified, all tasks will share the same sheet
) => {
  const databaseList = databaseUIDList.map((uid) =>
    useDatabaseV1Store().getDatabaseByUID(uid)
  );

  const databaseListGroupByEnvironment = groupBy(
    databaseList,
    (db) => db.effectiveEnvironment
  );
  const stageList = orderBy(
    Object.keys(databaseListGroupByEnvironment).map((env) => {
      const environment = useEnvironmentV1Store().getEnvironmentByName(env);
      const databases = databaseListGroupByEnvironment[env];
      return {
        environment,
        databases,
      };
    }),
    [(stage) => stage.environment?.order],
    ["asc"]
  );

  const steps: Plan_Step[] = [];
  for (let i = 0; i < stageList.length; i++) {
    const step = Plan_Step.fromJSON({});
    const { databases } = stageList[i];
    for (let j = 0; j < databases.length; j++) {
      const db = databases[j];
      const sqlIndex = databaseUIDList.findIndex((uid) => uid === db.uid);
      const spec = await buildSpecForTarget(db.name, params, sheetUID);
      step.specs.push(spec);
      maybeSetInitialSQLForSpec(spec, sqlIndex, params);
      maybeSetInitialDatabaseConfigForSpec(spec, params);
    }
    steps.push(step);
  }
  return steps;
};

export const buildStepsForDatabaseGroup = async (
  params: CreateIssueParams,
  databaseGroupName: string
) => {
  // Create sheet from SQL template in URL query
  // The sheet will be used when previewing rollout
  const sql = params.initialSQL.sql ?? "";
  const sheetCreate = {
    ...createEmptyLocalSheet(),
  };
  setSheetStatement(sheetCreate, sql);
  const sheet = await useSheetV1Store().createSheet(
    params.project.name,
    sheetCreate
  );

  const spec = await buildSpecForTarget(
    databaseGroupName,
    params,
    extractSheetUID(sheet.name)
  );
  const step = Plan_Step.fromJSON({
    specs: [spec],
  });
  return [step];
};

export const buildStepsViaDeploymentConfig = async (
  params: CreateIssueParams,
  sheetUID: string
) => {
  const { project } = params;
  const spec = await buildSpecForTarget(
    `${project.name}/deploymentConfigs/default`,
    params,
    sheetUID
  );
  maybeSetInitialSQLForSpec(spec, 0, params);
  maybeSetInitialDatabaseConfigForSpec(spec, params);
  const step = Plan_Step.fromPartial({
    specs: [spec],
  });
  return [step];
};

export const buildStepsViaChangelist = async (
  databaseUIDList: string[],
  changelistResourceName: string,
  params: CreateIssueParams
) => {
  const changelist = await useChangelistStore().getOrFetchChangelistByName(
    changelistResourceName
  );
  const { changes } = changelist;

  const databaseList = databaseUIDList.map((uid) =>
    useDatabaseV1Store().getDatabaseByUID(uid)
  );

  const databaseListGroupByEnvironment = groupBy(
    databaseList,
    (db) => db.effectiveEnvironment
  );
  const stageList = orderBy(
    Object.keys(databaseListGroupByEnvironment).map((env) => {
      const environment = useEnvironmentV1Store().getEnvironmentByName(env);
      const databases = databaseListGroupByEnvironment[env];
      return {
        environment,
        databases,
      };
    }),
    [(stage) => stage.environment?.order],
    ["asc"]
  );

  const steps: Plan_Step[] = [];
  for (let i = 0; i < stageList.length; i++) {
    const step = Plan_Step.fromJSON({});
    const { databases } = stageList[i];
    for (let j = 0; j < databases.length; j++) {
      const db = databases[j];
      for (let k = 0; k < changes.length; k++) {
        const change = changes[k];
        const statement = await generateSQLForChangeToDatabase(change, db);
        const sheetUID = nextUID();
        const sheetName = `${params.project.name}/sheets/${sheetUID}`;
        const sheet = getLocalSheetByName(sheetName);
        setSheetStatement(sheet, statement);
        const spec = await buildSpecForTarget(db.name, params, sheetUID);
        step.specs.push(spec);
      }
    }
    steps.push(step);
  }

  return steps;
};

export const buildSpecForTarget = async (
  target: string,
  { project, query }: CreateIssueParams,
  sheetUID?: string
) => {
  const sheet = `${project.name}/sheets/${sheetUID ?? nextUID()}`;
  const template = query.template as TemplateType | undefined;
  const spec = Plan_Spec.fromJSON({
    id: uuidv4(),
  });
  if (template === "bb.issue.database.data.update") {
    spec.changeDatabaseConfig = Plan_ChangeDatabaseConfig.fromJSON({
      target,
      type: Plan_ChangeDatabaseConfig_Type.DATA,
    });
    if (query.sheetId) {
      const sheet = await useSheetV1Store().getOrFetchSheetByUID(query.sheetId);
      if (sheet) {
        spec.changeDatabaseConfig.sheet = sheet.name;
      }
    }
    if (!spec.changeDatabaseConfig.sheet) {
      spec.changeDatabaseConfig.sheet = sheet;
    }
  }
  if (template === "bb.issue.database.schema.update") {
    const type =
      query.ghost === "1"
        ? Plan_ChangeDatabaseConfig_Type.MIGRATE_GHOST
        : Plan_ChangeDatabaseConfig_Type.MIGRATE;
    spec.changeDatabaseConfig = Plan_ChangeDatabaseConfig.fromJSON({
      target,
      type,
      sheet,
    });
  }
  if (template === "bb.issue.database.schema.baseline") {
    spec.changeDatabaseConfig = Plan_ChangeDatabaseConfig.fromJSON({
      target,
      type: Plan_ChangeDatabaseConfig_Type.BASELINE,
    });
  }

  return spec;
};

export const previewPlan = async (plan: Plan, params: CreateIssueParams) => {
  const rollout = await rolloutServiceClient.previewRollout({
    project: params.project.name,
    plan,
  });
  // Touch UIDs for each object for local referencing
  rollout.plan = plan.name;
  rollout.uid = nextUID();
  rollout.name = `${params.project.name}/rollouts/${rollout.uid}`;
  rollout.stages.forEach((stage) => {
    stage.uid = nextUID();
    stage.name = `${rollout.name}/stages/${stage.uid}`;
    stage.tasks.forEach((task) => {
      task.uid = nextUID();
      task.name = `${stage.name}/tasks/${task.uid}`;
    });
  });

  await maybeWrapStatementsAsSheets(plan, rollout, params);

  return reactive(rollout);
};

const maybeWrapStatementsAsSheets = (
  plan: Plan,
  rollout: Rollout,
  params: CreateIssueParams
) => {
  const { query } = params;
  if (!query.databaseGroupName) {
    return;
  }
  const { stages } = rollout;
  for (let i = 0; i < stages.length; i++) {
    const stage = stages[i];
    const { tasks } = stage;
    for (let j = 0; j < tasks.length; j++) {
      const task = tasks[j];
      // For database group changes, tasks returned by previewRollout
      // have `sheet` fields actually generated SQL statements instead
      // of sheet names.
      // So we create a local sheet for each task and set the statement
      // to the local task.
      const statement = sheetNameOfTaskV1(task);
      const sheetName = `${params.project.name}/sheets/${nextUID()}`;
      const sheet = getLocalSheetByName(sheetName);
      sheet.database = task.target;
      setSheetStatement(sheet, statement);
      setSheetNameForTask(task, sheetName);
    }
  }
};

const maybeSetInitialSQLForSpec = (
  spec: Plan_Spec,
  index: number,
  params: CreateIssueParams
) => {
  const sheet = sheetNameForSpec(spec);
  if (!sheet) return;
  const uid = extractSheetUID(sheet);
  if (!uid.startsWith("-")) {
    // If the sheet is a remote sheet, ignore initial SQL in URL query
    return;
  }
  // Priority: sqlList[index] -> sql -> nothing
  const sql = params.initialSQL.sqlList?.[index] ?? params.initialSQL.sql ?? "";
  if (sql) {
    const sheetEntity = getLocalSheetByName(sheet);
    setSheetStatement(sheetEntity, sql);
  }
};

const maybeSetInitialDatabaseConfigForSpec = (
  spec: Plan_Spec,
  params: CreateIssueParams
) => {
  const branch = params.branch;
  if (!branch) {
    return;
  }

  const sheetName = sheetNameForSpec(spec);
  if (!sheetName) return;
  const uid = extractSheetUID(sheetName);
  if (!uid.startsWith("-")) {
    // If the sheet is a remote sheet, ignore initial Database configs.
    return;
  }

  const [_, sheetId] = getProjectAndSchemaDesignSheetId(branch);
  const temp = useSheetV1Store().getSheetByUID(sheetId);
  if (temp && temp.payload) {
    const sheetEntity = getLocalSheetByName(sheetName);
    sheetEntity.payload = temp.payload;
  }
};

const extractInitialSQLListFromQuery = (
  query: Record<string, string>
): {
  sqlList?: string[];
  sql?: string;
} => {
  const sqlListJSON = query.sqlList;
  if (sqlListJSON && sqlListJSON.startsWith("[") && sqlListJSON.endsWith("]")) {
    try {
      const sqlList = JSON.parse(sqlListJSON) as string[];
      if (Array.isArray(sqlList)) {
        if (sqlList.every((maybeSQL) => typeof maybeSQL === "string")) {
          return {
            sqlList,
          };
        }
      }
    } catch {
      // Nothing
    }
  }
  const sql = query.sql;
  if (sql && typeof sql === "string") {
    return {
      sql,
    };
  }
  return {};
};

export const prepareDatabaseList = async (
  databaseUIDList: string[],
  projectUID: string
) => {
  const databaseStore = useDatabaseV1Store();
  if (projectUID && projectUID !== String(UNKNOWN_ID)) {
    // For preparing the database if user visits creating issue url directly.
    // It's horrible to fetchDatabaseByUID one-by-one when query.databaseList
    // is big (100+ sometimes)
    // So we are fetching databaseList by project since that's better cached.
    const project = await useProjectV1Store().getOrFetchProjectByUID(
      projectUID
    );
    await prepareDatabaseListByProject(project.name);
  } else {
    // Otherwise, we don't have the projectUID (very rare to see, theoretically)
    // so we need to fetch the first database in databaseList by id,
    // and see what project it belongs.
    if (databaseUIDList.length > 0) {
      const firstDB = await databaseStore.getOrFetchDatabaseByUID(
        databaseUIDList[0]
      );
      if (databaseUIDList.length > 1) {
        await prepareDatabaseListByProject(firstDB.project);
      }
    }
  }
};

const prepareDatabaseListByProject = async (project: string) => {
  await useDatabaseV1Store().searchDatabaseList({
    parent: `instances/-`,
    filter: `project == "${project}"`,
  });
};

export const isValidStage = (stage: Stage): boolean => {
  for (const task of stage.tasks) {
    if (task.type === Task_Type.DATABASE_SCHEMA_BASELINE) {
      continue;
    }

    if (TaskTypeListWithStatement.includes(task.type)) {
      const sheetName = sheetNameOfTaskV1(task);
      const uid = extractSheetUID(sheetName);
      if (uid.startsWith("-")) {
        const sheet = getLocalSheetByName(sheetName);
        if (getSheetStatement(sheet).length === 0) {
          return false;
        }
      } else {
        const sheet = useSheetV1Store().getSheetByName(sheetName);
        if (!sheet) {
          return false;
        }
        if (getSheetStatement(sheet).length === 0) {
          return false;
        }
      }
    }
  }
  return true;
};
