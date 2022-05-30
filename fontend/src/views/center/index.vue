<template>
  <div class="container">
    <el-row align="middle" :gutter="24">
      <el-col :span="5">
        <el-avatar size="large" :src="assetsUrl + avatarUrl" />
      </el-col>
      <el-col :span="17">
        <p>{{author}}</p>
        <p>{{likeSum}}</p>
      </el-col>
    </el-row>

    <div class="info">
      <div class="list-flex">
        <div>头像</div>
        <el-avatar class="right-flex" size="small" :src="assetsUrl + avatarUrl" />
        <el-icon class="edit-btn"><ArrowRightBold @click="handleShow('avatar')" /></el-icon>
      </div>
      <div class="list-flex">
        <div>作者</div>
        <div class="right-flex">{{author}}</div>
        <el-icon class="edit-btn"><ArrowRightBold  @click="handleShow('author')" /></el-icon>
      </div>
      <div class="list-flex">
        <div>手机号</div>
        <div class="right-flex">{{phone}}</div>
        <el-icon class="edit-btn"><ArrowRightBold  @click="handleShow('phone')" /></el-icon>
      </div>
      <div class="list-flex">
        <div>邀请码</div>
        <div>{{inviteCode}}</div>
      </div>
      <div v-if="author !== ''" class="list-flex">
        <div>公钥</div>
        <div>
          <el-button size="small" type="success" :icon="Download" circle @click="downloadSecret(1)" />
        </div>
      </div>
      <div v-if="author !== ''" class="list-flex">
        <div>私钥</div>
        <div>
          <el-button size="small" type="success" :icon="Download" circle @click="downloadSecret(2)" />
        </div>
      </div>
    </div>

    <div class="bottom-flex">
      <el-button @click="singOut" class="sign-out" type="primary">退出登录</el-button>
    </div>

    <el-dialog
        v-model="dialogVisible"
        title="更新个人信息"
        top="25vh"
        width="70%"
        destroy-on-close
        :before-close="handleClose"
    >
      <div>
        <div v-if="infoType === 'avatar'">
          <el-upload
              class="upload-demo"
              :auto-upload="false"
              :limit="1"
              :file-list="avatarImg"
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
        </div>
        <div v-else>
          <div class="edit-flex">
            <div>{{infoType === 'author' ? '作者：' : '手机号：' }}</div>
            <el-input v-model="updateInput" class="input-style" placeholder="请输入"/>
          </div>
        </div>
      </div>
      <template #footer>
      <span class="dialog-footer">
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="updateAvatar()"
        >确认</el-button
        >
      </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
  import {ref} from "vue";
  import { useRouter } from 'vue-router';
  import type { UploadProps, UploadUserFile } from 'element-plus'
  import { Plus ,ArrowRightBold, Download } from '@element-plus/icons-vue'
  import avatar from '@/assets/img/avatar.jpeg';
  import userStore from "@/store/user";
  import routerStore from "@/store/router";
  import service from '@/api/http';

  const router = useRouter();

  const assetsUrl = ref(`http://view3d.${window.location.host}/`);

  const avatarUrl = ref('');

  const imageUrl = ref('');

  const avatarImg = ref<UploadUserFile[]>([]) // logo图片


  const likeSum = ref();
  const author = ref('');
  const phone = ref('');
  const inviteCode = ref('');
  const dialogVisible = ref(false)

  const infoType = ref('');
  const updateInput = ref('');

  // 页面初始化加载
  onMounted(() => {
    // 初始化根据当前地址栏pathname更新routerName
    author.value = userStore().getUserName
    phone.value = userStore().getPhone
    inviteCode.value = userStore().getInviteCode
    avatarUrl.value = userStore().getAvatarUrl
  })

  /**
   * 下载文件
   * @param content
   */
  const funDownload = function (content: string, type: number) {
    // 创建隐藏的可下载链接
    const eleLink = document.createElement('a');
    eleLink.download = type === 1 ? 'public.pem' : 'secretKey.pem';
    eleLink.style.display = 'none';
    // 字符内容转变成blob地址
    const blob = new Blob([content]);
    eleLink.href = URL.createObjectURL(blob);
    // 触发点击
    document.body.appendChild(eleLink);
    eleLink.click();
    // 然后移除
    document.body.removeChild(eleLink);
  };

  /**
   * 下载密钥
   */
  const downloadSecret = async (type: number) => {
    const result = await service.get('/user/secret-download', { username: userStore().getUserName, secrete_type: type, token: userStore().getToken })
    console.log(result);
    if (result) {
      funDownload(result, type)

      ElMessage({
        type: 'success',
        offset: 20,
        message: type === 1 ? '下载公钥成功！' : '下载密钥成功！'
      })
    }
  }

  /**
   * 打开弹窗
   */
  const handleShow = (type: string) => {
    dialogVisible.value = true;
    infoType.value = type;
  }

  /**
   * 关闭弹窗
   * @param done
   */
  const handleClose = (done: () => void) => {
    updateInput.value = '';
    done()
  }

  /**
   * 上传移除
   * @param file
   * @param uploadFiles
   */
  const handleRemove: UploadProps['onRemove'] = (file, uploadFiles) => {
    console.log(file, uploadFiles)
  }

  /**
   * 图片预览
   * @param uploadFile
   */
  const handlePreview: UploadProps['onPreview'] = (uploadFile) => {
    console.log(uploadFile)
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
      imageUrl.value = result.data?.file_addr;
    }
  }

  /**
   * 更新头像信息
   */
  const updateAvatar = async () => {
    const result = await service.post('/user/avatar', { token: userStore().getToken, username: userStore().getUserName, avatarUrl: imageUrl.value })
    if (result?.data) {
      avatarUrl.value = result?.data?.AvatarUrl || '';
      dialogVisible.value = false;

      ElMessage({
        type: 'success',
        offset: 20,
        message: '更新头像成功！',
      });
    }
  }

  /**
   * 退出接口
   */
  const singOut = async () => {
    return false;
  }

</script>

<style lang="scss" scoped>
  .container {
    position: relative;
    width: 100vw;
    height: 100vh;
    background-color: #d9d0f3;
    padding: 1.5rem;

    .info {
      color: #fff;
      width: 100%;
      margin-top: 1rem;
      padding: 0.3rem 0.8rem;
      border-radius: 6px;
      background-color: #544199;

      .list-flex {
        position: relative;
        display: flex;
        flex-flow: row wrap;
        justify-content: space-between;
        align-items: center;
        height: 3rem;
        margin: 0.5rem 0;
        border-bottom: 1px solid #fff;

        .right-flex {
          margin-right: 1.2rem;
        }

        .edit-btn {
          position: absolute;
          right: 0;
          top: 50%;
          transform: translateY(-50%);
          color: #fff;
        }
      }
    }

    .bottom-flex {
      display: flex;
      justify-content: center;
      margin-top: 5rem;

      .sign-out {
        width: 10rem;
        height: 2.5rem;
        background-color: #544199;
      }
    }

    .edit-flex {
      display: flex;
      align-items: center;

      .input-style {
        width: 10rem;
      }
    }
  }
</style>
