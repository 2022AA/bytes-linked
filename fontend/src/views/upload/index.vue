<template>
  <div class="container">
    <el-form
        ref="ruleFormRef"
        :model="ruleForm"
        :rules="rules"
        class="demo-ruleForm"
    >
      <el-form-item label="模型名称" prop="modelName">
        <el-input v-model="ruleForm.modelName" placeholder="请输入"/>
      </el-form-item>
      <el-form-item label="上传文件" prop="modelFile">
        <el-upload
            ref="uploadRef"
            class="upload-demo"
            :auto-upload="false"
            :before-upload="beforeFilesUpload"
            :on-change="onFileChange"
        >
          <template #trigger>
            <el-button type="primary">选择文件</el-button>
          </template>

          <template #tip>
            <div class="el-upload__tip">
              请上传压缩包格式
            </div>
          </template>
        </el-upload>
      </el-form-item>
      <el-form-item label="模型封面" prop="modelLogo">
        <el-upload
            class="upload-demo"
            :auto-upload="false"
            :limit="1"
            :file-list="modelLogo"
            :on-preview="handlePreview"
            :on-remove="handleRemove"
            :on-change="onImageChange"
            list-type="picture"
        >
          <el-button type="primary">选择文件</el-button>
          <template #tip>
            <div class="el-upload__tip">
              jpg/png 格式
            </div>
          </template>
        </el-upload>
      </el-form-item>
      <el-form-item label="上传密钥" prop="secretKey">
        <el-upload
            ref="secretKey"
            class="upload-demo"
            :auto-upload="false"
            :before-upload="beforeFilesUpload"
            :on-change="onPemFileChange"
        >
          <template #trigger>
            <el-button type="primary">选择文件</el-button>
          </template>

          <template #tip>
            <div class="el-upload__tip">
              请上传pem格式
            </div>
          </template>
        </el-upload>
      </el-form-item>
    </el-form>

    <div class="bottom-flex">
      <el-button class="upload" type="primary" @click="submitForm(ruleFormRef)">上传</el-button>
    </div>
  </div>
</template>

<script lang="ts" setup>
/* eslint-disable */
import {useRouter} from 'vue-router';
import type {FormInstance, FormRules, UploadInstance, UploadProps, UploadUserFile, UploadRawFile} from 'element-plus'
import routerStore from '@/store/router';
import service from '@/api/http';
import { resultType } from '@/types';
import userStore from "@/store/user";

const router = useRouter();

// 上传组装数据
const ruleForm = reactive({
  modelName: '',
  modelFile: '',
  modelLogo: '',
  modelUrl: '',
  basePath: '',
});

const uploadRef = ref<UploadInstance>(); // 模型文件
const secretKey = ref<UploadInstance>(); // 模型文件
const modelLogo = ref<UploadUserFile[]>([]) // 模型封面
const ruleFormRef = ref<FormInstance>(); // 模型规则
const pemKey = ref(''); // 密钥

/**
 * 表单规则
 */
const rules = reactive<FormRules>({
  modelName: [
    {required: true, message: '请输入', trigger: 'blur'},
    {min: 1, max: 20, message: '超过最长10个字符限制', trigger: 'blur'},
  ],
  modelFile: [
    {required: false, message: '请上传模型文件', trigger: 'blur'},
  ],
  modelLogo: [
    {required: false, message: '请上传模型封面', trigger: 'blur'},
  ],
})

const handleRemove: UploadProps['onRemove'] = (uploadFile, uploadFiles) => {
  console.log(uploadFile, uploadFiles)
}

const handlePreview: UploadProps['onPreview'] = (file) => {
  console.log(file)
}

const beforeFilesUpload: UploadProps['beforeUpload'] = (file) => {
  // console.log(file, 1)
}

/**
 * 文件Change事件
 * @param file
 * @param files
 */
const onFileChange: UploadProps['onChange'] = async (file, files) => {
  const formData = new FormData();
  formData.append('model', file.raw || '');
  const result = await service.upload('/file/upload', formData);
  if (result?.data) {
    ruleForm.modelFile = result.data?.file_addr;
    ruleForm.modelUrl = result.data?.file_addr;
    ruleForm.basePath = result.data?.file_sha1;
  }
  console.log("上传模型结果", result)
}

/**
 * 上传密钥文件
 * @param file
 * @param files
 */
const onPemFileChange: UploadProps['onChange'] = async (file, files) => {
  const formData = new FormData();
  formData.append('key', file.raw || '');
  const result = await service.upload('/file/upload', formData)
  if (result?.data) {
    pemKey.value = result.data;
  }
  console.log("上传密钥结果", result)
}

/**
 * 图片Change事件
 * @param file
 * @param files
 */
const onImageChange: UploadProps['onChange'] = async (file, files) => {
  const formData = new FormData();
  formData.append('image', file.raw || '');
  const result = await service.upload('/file/upload', formData)
  if (result?.data) {
    ruleForm.modelLogo = result.data?.file_addr;
  }
  console.log("上传图片结果", result)
}

/**
 * 上传模型
 * @param submitData
 */
const commitModelFn = async (submitData: object) => {
  const result = await service.post('/file/commit', submitData);
  console.log(result);

  if (result?.code === 0) {

    ElMessage({
      type: 'success',
      offset: 20,
      message: '上传作品成功！'
    })

    setTimeout(() => {
      // 跳转路由页
      router.push(`/collect`);
    }, 500);
  }
}

/**
 * 提交上传数据
 * @param formEl
 */
const submitForm = async (formEl: FormInstance | undefined) => {
  if (!formEl) return
  await formEl.validate((valid, fields) => {
    if (valid) {
      console.log(ruleForm);

      const submitData = {
        token: userStore().getToken,
        uid: Number(userStore().getUserUid),
        file_name: ruleForm.modelName,
        file_sha1: ruleForm.basePath,
        file_addr: ruleForm.modelUrl,
        img_addr: ruleForm.modelLogo,
        fisco_key: pemKey.value,
      }

      commitModelFn(submitData);

      console.log('submit!')
    } else {
      console.log('error submit!', fields)
    }
  })
}
</script>

<style lang="scss" scoped>
.container {
  position: relative;
  width: 100vw;
  height: 100vh;
  background-color: #d9d0f3;
  padding: 1.5rem;

  .bottom-flex {
    display: flex;
    justify-content: center;
    margin-top: 5rem;

    .upload {
      width: 10rem;
      height: 2.5rem;
      border-radius: 10px;
      background-color: #544199;
    }
  }
}
</style>
