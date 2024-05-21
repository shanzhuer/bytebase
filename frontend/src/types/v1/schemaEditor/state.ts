import type { ComposedDatabase } from "..";
import type { Branch } from "../../proto/v1/branch_service";
import type { Schema } from "./atomType";

export enum SchemaEditorTabType {
  TabForDatabase = "database",
  TabForTable = "table",
}

// Tab context for editing database.
export interface DatabaseTabContext {
  id: string;
  type: SchemaEditorTabType.TabForDatabase;
  // Parent could be a database.
  parentName: string;
  selectedSchemaId?: string;
}

// Tab context for editing table.
export interface TableTabContext {
  id: string;
  type: SchemaEditorTabType.TabForTable;
  // Parent could be a database or a branch.
  parentName: string;
  schemaId: string;
  tableId: string;
}

export type TabContext = {
  name?: string;
} & (DatabaseTabContext | TableTabContext);

export interface DatabaseSchema {
  database: ComposedDatabase;
  schemaList: Schema[];
  originSchemaList: Schema[];
}

export interface BranchSchema {
  branch: Branch;
  schemaList: Schema[];
  originSchemaList: Schema[];
}
