<template>
  <div class="container-item">
    <div v-show="activeName === '1'" class="upload" @click="uploadFiles()">
      <el-icon color="#ffffff" :size="30"><UploadFilled /></el-icon>
    </div>
    <el-tabs v-model="activeName" class="top-tabs">
      <el-tab-pane label="我的作品" name="1">
        <template v-for="(item, index) in modelList" :key="index">
          <div v-if="item.img_addr !== ''" class="list-row-flex">
            <div class="img-flex">
              <el-image class="list-img" :src=" assetsUrl + item.img_addr || ''" @click="toModelViewer(item.file_addr, item.file_sha1)" />
            </div>
            <div class="list-row">
              <div>{{item.file_name || ''}}</div>
              <div class="right-flex">
                <el-image v-show="canIntoAr" @click="toArShow(item.file_sha1)" class="ar-style" :src="arImg" />
                <el-icon @click="shareModel(item.file_addr, item.file_sha1)" color="#ffffff" :size="20" class="share-model"><Share /></el-icon>
                <el-icon @click="postModel(item.file_id)" color="#ffffff" :size="20"><Position /></el-icon>
              </div>
            </div>
            <div class="list-row">
              <div class="like">
<!--                <el-button-->
<!--                    type="warning"-->
<!--                    :icon="Star"-->
<!--                    size="small"-->
<!--                    circle-->
<!--                    :disabled="!item.can_star"-->
<!--                    @click="clickLike(item.file_id)"-->
<!--                />-->
                <p>{{item.like_cnt}}币</p>
              </div>
              <div>
                <el-button size="small" @click="confirmProve(item.file_id)">确权证明</el-button>
                <el-button size="small" @click="showDialog(item.file_id)">转让</el-button>
              </div>
            </div>
          </div>
        </template>
      </el-tab-pane>
      <el-tab-pane label="交易记录" name="2">
        <template v-for="(item, index) in recordList" :key="index">
          <div class="list-column-flex">
            <div class="img-flex">
              <el-image class="list-img" :src="item.img" />
            </div>
            <div class="list-column">
              <div>{{item.title}}</div>
              <div>原持有者：{{item.originalOwner}}</div>
              <div>现持有者：{{item.currentOwner}}</div>
              <div>交易时间：{{item.tradeTime}}</div>
            </div>
            <span>{{item.status}}</span>
          </div>
        </template>
      </el-tab-pane>
    </el-tabs>

    <el-dialog
        v-model="qrCodeVisible"
        title="作品二维码"
        top="25vh"
        width="75%"
        destroy-on-close
        :before-close="qrCodeClose"
    >
      <div style="text-align: center;">
        <vue-qrcode
            :value="qrCodeUrl"
            :options="{
            color: {
              dark: '#222222',
              light: '#fff',
            },
          }"
        ></vue-qrcode>
      </div>
    </el-dialog>

    <el-dialog
        v-model="dialogVisible"
        title="作品转让"
        width="80%"
        top="30vh"
        destroy-on-close
        :before-close="handleClose"
    >
      <div>
        <el-form
            ref="ruleFormRef"
            :model="ruleForm"
            :rules="rules"
            class="demo-ruleForm"
        >
          <el-form-item label="对方用户名" prop="toUsername">
            <el-input placeholder="请输入" v-model="ruleForm.toUsername" />
          </el-form-item>
          <el-form-item label="密钥" prop="secretKey">
            <el-upload
                ref="secretKey"
                class="upload-demo"
                :auto-upload="false"
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
      </div>
      <template #footer>
      <span class="dialog-footer">
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmTrans()">确认转出</el-button>
      </span>
      </template>
    </el-dialog>

    <el-dialog
        v-model="proveVisible"
        title="确权证明"
        top="25vh"
        width="70%"
        destroy-on-close
        :before-close="qrCodeClose"
    >
      <div>
        证明ID：{{mySha1}}
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
  import {ref} from "vue";
  import { useRouter } from 'vue-router';
  import type { FormInstance, FormRules, UploadInstance, UploadProps, UploadUserFile } from 'element-plus'
  import { UploadFilled, Search, Share, Position, Camera, Star } from '@element-plus/icons-vue'
  import demo from '@/assets/img/img-demo.png'
  import arImg from '@/assets/img/ar-img.png'
  import routerStore from '@/store/router';
  import userStore from '@/store/user';
  import service from '@/api/http';
  import { itemType } from '@/types'

  const router = useRouter();

  const assetsUrl = ref(`http://view3d.${window.location.host}/`);

  const qrCodeUrl = ref('http://2022aa.2022AA.com/#/');

  const activeName = ref('1'); // 激活tabs
  const dialogVisible = ref(false)
  const qrCodeVisible = ref(false);
  const proveVisible = ref(false);
  const canIntoAr = ref(false); // 是否能进入AR模式

  // 上传组装数据
  const ruleForm = reactive({
    toUsername: '',
    privateKey: '',
  });

  const currentId = ref(0);

  const mySha1 = ref('');

  const uploadRef = ref<UploadInstance>(); // 模型文件

  const modelList = reactive([
    {
      file_id: 0,
      file_addr: '',
      file_sha1: '',
      img_addr: '',
      file_name: '',
      like_cnt: 0,
      status: 0,
      can_star: true,
    },
  ])

  const recordList = reactive([
    { img: demo, title: '作品一', originalOwner: 'CkarkChu1', currentOwner: '2022AA2', tradeTime: '2022-02-02 19:00', status: '成功' },
  ])

  onMounted(async () => {
    // 判断是否支持AR识别
    const userAgentInfo = window.navigator.userAgent;
    if (/iPhone/.test(userAgentInfo) && /Safari/.test(userAgentInfo)) {
      canIntoAr.value = true;
    }

    // 获取我的作品列表
    const result = await service.get('/file/list_own', { uid: userStore().getUserUid })
    if (result?.code === 0) {
      const list = result.data;
      if (list?.length) {
        list.forEach((item: itemType, index: number) => {
          if (modelList[index] === undefined) {
            modelList[index] = {
              file_id: 0,
              file_addr: '',
              file_sha1: '',
              img_addr: '',
              file_name: '',
              like_cnt: 0,
              status: 1,
              can_star: true,
            }
          }
          modelList[index].file_id = item?.file_id || 0;
          modelList[index].file_addr = item?.file_addr || '';
          modelList[index].file_sha1 = item?.file_sha1 || '';
          modelList[index].img_addr = item?.img_addr || '';
          modelList[index].file_name = item?.file_name || '';
          modelList[index].like_cnt = item?.like_cnt || 0;
          modelList[index].status = item?.status || 0;
          modelList[index].can_star = true;
        })
      }
    }
  })

  /**
   * 表单规则
   */
  const rules = reactive<FormRules>({
    toUsername: [
      {required: true, message: '请输入', trigger: 'blur'},
      {min: 1, max: 20, message: '超过最长10个字符限制', trigger: 'blur'},
    ],
  })

  /**
   * 提交上传数据
   * @param formEl
   */
  const submitForm = async (formEl: FormInstance | undefined) => {
    if (!formEl) return
    await formEl.validate((valid, fields) => {
      if (valid) {
        console.log('submit!')
      } else {
        console.log('error submit!', fields)
      }
    })
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
      ruleForm.privateKey = result.data;
    }
    console.log("上传密钥结果", result)
  }

  const showDialog = (fileId: number) => {
    dialogVisible.value = true;
    currentId.value = fileId;
  }

  /**
   * 转出作品
   * @param fileId
   */
  const confirmTrans = async () => {
    const submitData = {
      token: userStore().getToken,
      uid: Number(userStore().getUserUid),
      transferFileId: currentId.value,
      toUsername: ruleForm.toUsername,
      privateKey: ruleForm.privateKey,
    }
    const result = await service.post('/user/transaction', submitData);
    console.log(result);
    if (result?.code === 0) {
      dialogVisible.value = false;

      ElMessage({
        type: 'success',
        offset: 20,
        message: '作品转出成功！',
      });

      setTimeout(() => {
        router.push('/collect');
      }, 500);
    }
  }

  /**
   * 确权证明
   * @param fileId
   */
  const confirmProve = async (fileId: number) => {
    const submitData = {
      token: userStore().getToken,
      uid: Number(userStore().getUserUid),
      fileId,
    }
    const result = await service.get('/user/file', submitData);
    console.log(result);
    if (result?.data) {
      mySha1.value = result.data.TransactionId;
      proveVisible.value = true;
    }
  }

  /**
   * 上传文件页面
   */
  const uploadFiles = () => {
    // 跳转路由页
    router.push(`/upload`);
  }

  /**
   * 关闭弹窗
   * @param done
   */
  const handleClose = (done: () => void) => {
    done()
  }

  /**
   * 关闭二维码弹窗
   * @param done
   */
  const qrCodeClose = (done: () => void) => {
    done()
  }

  /**
   * 发布我的作品
   * @param file_id
   */
  // eslint-disable-next-line camelcase
  const postModel = async (file_id: number) => {
    console.log('发布')
    // eslint-disable-next-line camelcase
    const result = await service.post('/file/public', { token: userStore().getToken, uid: Number(userStore().getUserUid), file_id, avatar_url: userStore().getAvatarUrl })
    console.log(result);
    if (result?.code === 0) {
      ElMessage({
        type: 'success',
        offset: 20,
        message: '作品发布成功！',
      });
    }
  }

  /**
   * 展示二维码弹窗
   * @param modelUrl
   * @param basePath
   */
  const shareModel = async (modelUrl: string, basePath: string) => {
    console.log('分享')
    qrCodeUrl.value = `${assetsUrl.value}?modelUrl=${encodeURIComponent(modelUrl)}&basePath=${encodeURIComponent(`assets/${basePath}/`)}`
    qrCodeVisible.value = true;
  }

  /**
   * 跳转模型详情页
   * @param modelUrl
   * @param basePath
   */
  const toModelViewer = (modelUrl: string, basePath: string) => {
    console.log(`${assetsUrl.value}?=modelUrl=${encodeURIComponent(modelUrl)}&basePath=${encodeURIComponent(basePath)}`)
    window.location.href = `${assetsUrl.value}?modelUrl=${encodeURIComponent(modelUrl)}&basePath=${encodeURIComponent(`assets/${basePath}/`)}`
  }

  /**
   * 跳转AR
   * @param sha1
   */
  const toArShow = (sha1: string) => {
    if (canIntoAr.value) {
      window.location.href = `${assetsUrl.value}assets/${sha1}/scene.usdz`
    } else {
      ElMessage({
        type: 'warning',
        offset: 20,
        message: '很抱歉，您的设备暂不支持AR深景模式查看，我们将尽快兼容。'
      })
    }
  }

  /**
   * 点赞作品
   * @param fileId
   */
  const clickLike = async (fileId: number) => {

    const submitData = {
      token: userStore().getToken,
      uid: Number(userStore().getUserUid),
      file_id: fileId,
      like_num: 1,
    }

    const result = await service.post('/file/like', submitData);
    console.log(result)
    if (result?.code === 0) {
      const currentList = Array.from(modelList)

      // 点赞成功暂时禁用点赞功能
      currentList.forEach((item, index) => {
        if (item.file_id === fileId) {
          modelList[index].can_star = false;
          modelList[index].like_cnt += 1;
        }
      })

      // 2秒后可重新点赞
      setTimeout(() => {
        currentList.forEach((item, index) => {
          if (item.file_id === fileId) {
            modelList[index].can_star = true;
          }
        })
      }, 1500)

    }
  }
</script>

<style lang="scss" scoped>
.container-item {
  width: 100vw;
  height: 100vh;
  overflow: auto;
  background-color: #d9d0f3;
  padding: 1rem 1rem 3.6rem;

  .upload {
    position: fixed;
    top: 1.2rem;
    right: 1rem;
    z-index: 100;
  }

  .list-row-flex {
    font-size: 0.9rem;
    width: 100%;
    color: #fff;
    margin-bottom: 0.8rem;
    background-color: #544199;
    padding-bottom: 0.6rem;
    border-radius: 6px;

    .share-model {
      margin: 0 0.5rem;
    }

    .ar-style {
      width: 2rem;
      height: 2rem;
    }

    .img-flex {
      width: 100%;
      height: auto;

      .list-img {
        width: 100%;
        height: 100%;
        border-top-left-radius: 6px;
        border-top-right-radius: 6px;
      }
    }

    .list-row {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 0.4rem 0.6rem;

      .right-flex {
        display: flex;
        align-items: center;
      }

      .like {
        display: flex;
        align-items: center;

        p {
          font-size: 1.1rem;
          margin-left: 0.4rem;
        }
      }
    }
  }

  .list-column-flex {
    position: relative;
    display: flex;
    align-items: center;
    font-size: 0.9rem;
    width: 100%;
    height: 6rem;
    color: #fff;
    margin-bottom: 0.8rem;
    background-color: #544199;
    border-radius: 6px;
    padding: 0 0.6rem;

    .img-flex {
      height: 5rem;
      width: 5rem;

      .list-img {
        width: 100%;
        height: 100%;
        border-top-left-radius: 6px;
        border-top-right-radius: 6px;
      }
    }

    .list-column {
      margin-left: 1rem;
      font-size: 0.8rem;
      color: #a2a3ab;

      &>:first-child {
        font-size: 1rem;
        color: #fff;
      }
    }

    span {
      position: absolute;
      right: 0.6rem;
      top: 0.6rem;
      color: #a2a3ab;
      font-size: 0.8rem;
    }
  }
}
</style>
