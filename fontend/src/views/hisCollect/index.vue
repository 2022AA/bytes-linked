<template>
  <div class="container">
    <el-row align="middle" :gutter="24">
      <el-col :span="5">
        <el-avatar size="large" :src=" assetsUrl + avatar || ''"/>
      </el-col>
      <el-col :span="17">
        <p class="author">{{author}}</p>
<!--        <el-icon color="#ffffff" :size="20"><ChatDotRound /></el-icon>-->
      </el-col>
    </el-row>

    <div class="list-title">
      他的作品（{{modelList?.length ?? 0}}）
    </div>

    <div>
      <template v-for="(item, index) in modelList" :key="index">
        <div v-if="item.img_addr !== ''" class="list-flex">
          <div class="img-flex">
            <el-image class="list-img" :src=" assetsUrl + item?.img_addr || ''" @click="toModelViewer(item.file_addr, item.file_sha1)" />
          </div>
          <div class="list-row">
            <div>{{item.file_name || '-'}}</div>
            <div class="right-flex">
              <el-image v-show="canIntoAr" @click="toArShow(item.file_sha1)" class="ar-style" :src="arImg" />
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
          </div>
        </div>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref } from "vue";
  import { ChatDotRound, Share, Position, Camera, Star } from '@element-plus/icons-vue'
  import { useRouter } from 'vue-router';
  import { itemType } from "@/types";
  import arImg from '@/assets/img/ar-img.png'
  import service from "@/api/http";
  import userStore from "@/store/user";

  const router = useRouter();

  const assetsUrl = ref(`http://view3d.${window.location.host}/`);
  const qrCodeUrl = ref('http://2022aa.2022AA.com/#/');

  const avatar = ref(''); // 头像
  const author = ref(''); // 用户名
  const canIntoAr = ref(false); // 是否能进入AR模式

  const modelList = reactive([
    {
      username: '',
      file_id: 0,
      file_addr: '',
      file_sha1: '',
      img_addr: '',
      file_name: '',
      like_cnt: 0,
      status: 0,
      avatar_url: '',
      can_star: true,
    },
  ])

  onMounted(async () => {
    // 获取他的uid
    const ownerUid = Number((window.location.hash).split('=')[1]);
    console.log('他的作品：', ownerUid);

    // 判断是否支持AR识别
    const userAgentInfo = window.navigator.userAgent;
    if (/iPhone/.test(userAgentInfo) && /Safari/.test(userAgentInfo)) {
      canIntoAr.value = true;
    }

    // 获取我的作品列表
    const result = await service.get('/file/list_own', { uid: ownerUid })
    if (result?.code === 0) {
      const list = result.data;
      if (list?.length) {
        avatar.value = list[0].avatar_url ?? '';
        author.value = list[0].username ?? '-';
        list.forEach((item: itemType, index: number) => {
          if (modelList[index] === undefined) {
            modelList[index] = {
              username: '',
              file_id: 0,
              file_addr: '',
              file_sha1: '',
              img_addr: '',
              file_name: '',
              like_cnt: 0,
              status: 1,
              avatar_url: '',
              can_star: true,
            }
          }

          modelList[index].username = item?.username || '';
          modelList[index].file_id = item?.file_id || 0;
          modelList[index].file_addr = item?.file_addr || '';
          modelList[index].file_sha1 = item?.file_sha1 || '';
          modelList[index].img_addr = item?.img_addr || '';
          modelList[index].file_name = item?.file_name || '';
          modelList[index].like_cnt = item?.like_cnt || 0;
          modelList[index].status = item?.status || 0;
          modelList[index].avatar_url = item?.avatar_url || '';
          modelList[index].can_star = true;
        })
      }
    }
  })

  const toNav = (path: string) => {
    // 跳转路由页
    router.push(path);
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
.container {
  width: 100vw;
  height: 100vh;
  overflow: auto;
  background-color: #d9d0f3;
  padding: 1rem 1rem 3.6rem;

  .author {
    color: #544199;
    margin-bottom: 0.4rem;
  }

  .list-title {
    margin: 1rem 0;
  }

  .list-flex {
    font-size: 0.9rem;
    width: 100%;
    color: #fff;
    margin-bottom: 1rem;
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

        .like {
          display: flex;
          align-items: center;
          margin-left: 0.5rem;

          p {
            font-size: 1.1rem;
            margin-left: 0.4rem;
          }
        }
      }
    }
  }
}
</style>
