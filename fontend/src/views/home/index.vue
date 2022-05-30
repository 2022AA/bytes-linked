<template>
  <div class="container-item">
<!--    <div class="slide">-->
<!--      <el-image class="nav-img" :src="nav"></el-image>-->
<!--    </div>-->

    <div class="search">
      <el-input
          v-model="keyWords"
          class="search-input"
          placeholder="请输入作品名称"
          style="{ 'border-radius': '20px' }"
          :prefix-icon="Search"
      />
<!--      <el-button class="search-btn" type="text">搜索</el-button>-->
    </div>

    <div>
      <template v-for="(item, index) in modelList" :key="index">
        <div v-if="item.img_addr !== ''" class="list-flex">
          <div class="list-desc">
            <div class="author" @click="toHisCollect(item?.owner_uid || 0)">
              <el-avatar size="small" :src=" assetsUrl + item.avatar_url || ''"/>
              <p>{{item.username}}</p>
            </div>
            <div class="icon-flex">
              <el-image v-show="canIntoAr" @click="toArShow(item.file_sha1)" class="ar-style" :src="arImg" />
              <el-icon @click="shareModel(item.file_addr, item.file_sha1)" color="#ffffff" :size="20" class="share-model"><Share /></el-icon>
            </div>
          </div>
          <div class="list-desc">
            <div>{{item.file_name}}</div>
            <div class="like">
              <el-button
                type="warning"
                :icon="Star"
                size="small"
                circle
                :disabled="!item.can_star"
                @click="clickLike(item.file_id)"
              />
              <p>{{item.like_cnt}}币</p>
            </div>
          </div>
          <div class="img-flex">
            <el-image class="list-img" :src=" assetsUrl + item.img_addr || ''" @click="toModelViewer(item.file_addr, item.file_sha1)" />
          </div>
        </div>
      </template>
    </div>

    <el-dialog
        v-model="dialogVisible"
        title="温馨提示"
        top="25vh"
        width="70%"
        destroy-on-close
        :before-close="handleClose"
    >
      <div class="home-tips">
        <p>如未下载密钥，请先下载<br/> 该密钥将作为个人身份重要凭证</p>
        <div>
          <el-button type='primary' @click="downloadSecret(1)">公钥下载</el-button>
          <el-button type="primary" @click="downloadSecret(2)">私钥下载</el-button>
        </div>
      </div>
    </el-dialog>

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
  </div>
</template>

<script setup lang="ts">
  import {ref} from "vue";
  import { useRouter } from 'vue-router';
  import { Search, Share, Camera, Medal, Star } from '@element-plus/icons-vue'
  import nav from '@/assets/img/nav.png'
  import demo from '@/assets/img/img-demo.png'
  import arImg from '@/assets/img/ar-img.png'
  import routerStore from '@/store/router';
  import userStore from '@/store/user';
  import service from '@/api/http';
  import {itemType} from "@/types";

  const router = useRouter();
  const { t } = useI18n();

  const assetsUrl = ref(`http://view3d.${window.location.host}/`);
  const keyWords = ref(''); // 搜索关键字
  const qrCodeUrl = ref('http://2022aa.2022AA.com/#/');
  const dialogVisible = ref(false);
  const qrCodeVisible = ref(false);
  const canIntoAr = ref(false); // 是否能进入AR模式

  const modelList = reactive([
    {
      file_id: 0,
      file_addr: '',
      file_sha1: '',
      avatar_url: '',
      img_addr: '',
      file_name: '',
      like_cnt: 0,
      username: '',
      ar_tag: false,
      owner_uid: 0,
      can_star: true,
    },
  ])

  onMounted(async () => {
    const routerFrom = routerStore().getRouterFrom;
    const routerName = routerStore().getRouterName;
    const storagePem = window.localStorage.getItem('downloadPem');

    // 如果是从注册页面过来的，提示下载密钥
    // if (routerFrom === '/' && routerName === '/home' && (storagePem == null || storagePem === '')) {
    //   dialogVisible.value = true;
    // }
    // console.log(`from：${routerFrom}`, `to：${routerName}`)

    // 判断是否支持AR识别
    const userAgentInfo = window.navigator.userAgent;
    if (/iPhone/.test(userAgentInfo) && /Safari/.test(userAgentInfo)) {
      canIntoAr.value = true;
    }

    const result = await service.get('/file/list_public')
    console.log(result);

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
              username: '',
              ar_tag: false,
              avatar_url: '',
              owner_uid: 0,
              can_star: true,
            }
          }
          modelList[index].file_id = item?.file_id || 0;
          modelList[index].file_addr = item?.file_addr || '';
          modelList[index].file_sha1 = item?.file_sha1 || '';
          modelList[index].img_addr = item?.img_addr || '';
          modelList[index].file_name = item?.file_name || '';
          modelList[index].like_cnt = item?.like_cnt || 0;
          modelList[index].username = item?.username || '';
          modelList[index].ar_tag = item?.ar_tag || false;
          modelList[index].avatar_url = item?.avatar_url || '';
          modelList[index].owner_uid = item?.owner_uid || 0;
          modelList[index].can_star = true;
        })
      }
    }
  })

  /**
   * 查看他的作品集合
   */
  const toHisCollect = (ownerUid: number) => {
    // 跳转路由页
    router.push(`/hisCollect?ownerUid=${ownerUid}`);
  }

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

      if (type === 2) {
        // 下载完成，关闭提示弹窗
        dialogVisible.value = false;

        // 记录该用户已下载私钥的行为
        window.localStorage.setItem('downloadPem','true');
      }

      ElMessage({
        type: 'success',
        offset: 20,
        message: type === 1 ? '下载公钥成功！' : '下载密钥成功！'
      })
    }
  }

  /**
   * 关闭密钥下载弹窗
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
    if (userStore().getBalance !== '') {
      const lastBalance = Number(userStore().getBalance);
      if (lastBalance < 1) {
        ElMessage({
          type: 'warning',
          offset: 20,
          message: '抱歉，您的剩余币不足，不可以点赞！',
        })
      } else {
        const submitData = {
          token: userStore().getToken,
          uid: Number(userStore().getUserUid),
          file_id: fileId,
          like_num: 1,
        }

        const result = await service.post('/file/like', submitData);
        console.log(result)
        if (result?.data) {

          // 更新剩余币
          userStore().$patch({
            balance: result.data.toString(),
          })
          window.localStorage.setItem('balance', (result.data).toString());

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
    } else {
      ElMessage({
        type: 'warning',
        offset: 20,
        message: '抱歉，您的剩余币不足，不可以点赞！',
      })
    }
  }
</script>

<style lang="scss" scoped>
  .container-item {
    width: 100vw;
    height: 100vh;
    overflow: auto;
    background-color: #d9d0f3;
    padding: 0.5rem 1rem 3.6rem;

    .share-model {
      margin-left: 0.5rem;
    }

    .home-tips {
      text-align: center;

      p {
        font-size: 0.95rem;
        color: #ff4500;
        margin-bottom: 1rem;
      }
    }

    .search {
      display: flex;
      justify-content: center;
      align-items: center;
      width: 100%;
      margin: 1rem auto;

      .search-input {
        width: 100%;
        height: 2.4rem;
      }

      .search-btn {
        font-size: 1.1rem;
      }
    }

    .slide{
      margin: 1rem 0;

      .nav-img {
        border-radius: 10px;
      }
    }

    .list-flex {
      width: 100%;
      border-radius: 6px;
      background-color: #544199;
      margin-bottom: 1rem;
      padding: 0.8rem 0.8rem 0.5rem;

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
          margin-top: 0.5rem;
          border-radius: 6px;
        }
      }

      .icon-flex {
        display: inline-flex;
        align-items: center;
      }

      .list-desc {
        position: relative;
        display: flex;
        justify-content: space-between;
        align-items: center;
        color: #fff;
        padding: 0.3rem 0;

        .author {
          display: flex;
          align-items: center;

          p {
            font-size: 1.1rem;
            color: #e6a23c;
            margin-left: 0.8rem;
            text-decoration: underline;
          }
        }

        .like {
          display: flex;
          align-items: center;

          p {
            font-size: 1.1rem;
            margin-left: 0.5rem;
          }
        }
      }
    }
  }
</style>
