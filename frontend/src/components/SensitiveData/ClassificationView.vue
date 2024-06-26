<template>
  <div class="w-full space-y-4">
    <div class="flex items-center justify-end">
      <NButton
        type="primary"
        :disabled="!hasPermission || !hasSensitiveDataFeature"
        @click="onUpload"
      >
        {{ $t("settings.sensitive-data.classification.upload") }}
      </NButton>
      <input
        ref="uploader"
        type="file"
        accept=".json"
        class="sr-only hidden"
        :disabled="!hasPermission || !hasSensitiveDataFeature"
        @input="onFileChange"
      />
    </div>
    <div class="textinfolabel">
      {{ $t("settings.sensitive-data.classification.upload-label") }}
      <span
        class="normal-link cursor-pointer hover:underline"
        @click="state.showExampleModal = true"
      >
        {{ $t("settings.sensitive-data.view-example") }}
      </span>
    </div>
    <div
      v-if="settingStore.classification.length === 0"
      class="flex justify-center border-2 border-gray-300 border-dashed rounded-md relative h-72"
    >
      <SingleFileSelector
        class="space-y-1 text-center flex flex-col justify-center items-center absolute top-0 bottom-0 left-0 right-0"
        :support-file-extensions="['.json']"
        :max-file-size-in-mi-b="maxFileSizeInMiB"
        :disabled="!hasPermission || !hasSensitiveDataFeature"
        @on-select="onFileSelect"
      >
        <template #image>
          <NoDataPlaceholder
            :border="false"
            :img-attrs="{ class: '!max-h-[10vh]' }"
          />
        </template>
      </SingleFileSelector>
    </div>
    <div v-else class="h-full">
      <ClassificationTree
        :classification-config="settingStore.classification[0]"
      />
    </div>
  </div>

  <DataExampleModal
    v-if="state.showExampleModal"
    :example="JSON.stringify(example, null, 2)"
    @dismiss="state.showExampleModal = false"
  />

  <BBAlert
    v-model:show="state.showOverrideModal"
    type="warning"
    :ok-text="$t('settings.sensitive-data.classification.override-confirm')"
    :title="$t('settings.sensitive-data.classification.override-title')"
    :description="$t('settings.sensitive-data.classification.override-desc')"
    @ok="upsertSetting"
    @cancel="state.showOverrideModal = false"
  />
</template>

<script lang="ts" setup>
import { v4 as uuidv4 } from "uuid";
import { computed, reactive, ref, watch } from "vue";
import { useI18n } from "vue-i18n";
import {
  featureToRef,
  useCurrentUserV1,
  useSettingV1Store,
  pushNotification,
} from "@/store";
import type {
  DataClassificationSetting_DataClassificationConfig_Level as ClassificationLevel,
  DataClassificationSetting_DataClassificationConfig_DataClassification as DataClassification,
} from "@/types/proto/v1/setting_service";
import { DataClassificationSetting_DataClassificationConfig } from "@/types/proto/v1/setting_service";
import { hasWorkspacePermissionV2 } from "@/utils";

const uploader = ref<HTMLInputElement | null>(null);
const maxFileSizeInMiB = 10;

interface UploadClassificationConfig {
  title: string;
  levels: ClassificationLevel[];
  classifications: DataClassification[];
}

interface LocalState {
  classification?: DataClassificationSetting_DataClassificationConfig;
  showExampleModal: boolean;
  showOverrideModal: boolean;
}

const { t } = useI18n();
const settingStore = useSettingV1Store();
const currentUser = useCurrentUserV1();
const state = reactive<LocalState>({
  showExampleModal: false,
  showOverrideModal: false,
});

watch(
  () => state.classification,
  async () => {
    if (settingStore.classification.length !== 0) {
      state.showOverrideModal = true;
      return;
    }

    await upsertSetting();
  }
);

const upsertSetting = async () => {
  if (!state.classification) {
    return;
  }
  await settingStore.upsertSetting({
    name: "bb.workspace.data-classification",
    value: {
      dataClassificationSettingValue: {
        configs: [state.classification],
      },
    },
  });
  pushNotification({
    module: "bytebase",
    style: "SUCCESS",
    title: t("settings.sensitive-data.classification.upload-succeed"),
  });
  state.showOverrideModal = false;
};

const hasPermission = computed(() => {
  return hasWorkspacePermissionV2(currentUser.value, "bb.policies.update");
});
const hasSensitiveDataFeature = featureToRef("bb.feature.sensitive-data");

const onUpload = () => {
  uploader.value?.click();
};

const onFileChange = () => {
  const files: File[] = (uploader.value as any).files;
  if (files.length !== 1) {
    return;
  }
  const file = files[0];
  if (file.size > maxFileSizeInMiB * 1024 * 1024) {
    pushNotification({
      module: "bytebase",
      style: "CRITICAL",
      title: t("common.file-selector.size-limit", {
        size: maxFileSizeInMiB,
      }),
    });
    return;
  }
  onFileSelect(file);
};

const onFileSelect = (file: File) => {
  const fr = new FileReader();
  fr.onload = () => {
    if (!fr.result) {
      return;
    }
    const data: UploadClassificationConfig = JSON.parse(fr.result as string);
    if (
      !Array.isArray(data.classifications) ||
      data.classifications.length === 0
    ) {
      return pushNotification({
        module: "bytebase",
        style: "CRITICAL",
        title: "Data format error",
        description: "Should has classifications array field",
      });
    }
    if (!Array.isArray(data.levels) || data.levels.length === 0) {
      return pushNotification({
        module: "bytebase",
        style: "CRITICAL",
        title: "Data format error",
        description: "Should has levels array field",
      });
    }
    if (
      data.classifications.length !==
      new Set(data.classifications.map((item) => item.id)).size
    ) {
      return pushNotification({
        module: "bytebase",
        style: "CRITICAL",
        title: "Data format error",
        description: "Should not contains duplicate classification id",
      });
    }
    state.classification =
      DataClassificationSetting_DataClassificationConfig.fromPartial({
        id: uuidv4(),
        title: data.title || "",
        levels: data.levels,
        classification: data.classifications.reduce(
          (obj, item) => {
            obj[item.id] = item;
            return obj;
          },
          {} as { [key: string]: DataClassification }
        ),
      });
  };
  fr.onerror = () => {
    pushNotification({
      module: "bytebase",
      style: "CRITICAL",
      title: "Read file error",
      description: String(fr.error),
    });
    return;
  };
  fr.readAsText(file);
};

const example: UploadClassificationConfig = {
  title: "Classification Example",
  levels: [
    {
      id: "1",
      title: "Level 1",
      description: "",
    },
    {
      id: "2",
      title: "Level 2",
      description: "",
    },
  ],
  classifications: [
    {
      id: "1",
      title: "Basic",
      description: "",
    },
    {
      id: "1-1",
      title: "Basic",
      description: "",
    },
    {
      id: "1-2",
      title: "Assert",
      description: "",
    },
    {
      id: "1-3",
      title: "Contact",
      description: "",
    },
    {
      id: "1-4",
      title: "Health",
      description: "",
    },
    {
      id: "2",
      title: "Relationship",
      description: "",
    },
    {
      id: "2-1",
      title: "Social",
      description: "",
    },
    {
      id: "2-2",
      title: "Business",
      description: "",
    },
  ],
};
</script>
