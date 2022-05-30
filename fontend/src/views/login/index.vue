<template>
  <div class="login-bg">
    <div class="login-body">
      <div class="top-title">
        <el-image :src="logoUrl"></el-image>
      </div>
      <el-form ref="ruleFormRef" :model="ruleForm" status-icon :rules="rules" class="form-style">
        <el-form-item class="list-style" prop="userName">
          <el-input size="large" v-model="ruleForm.userName" placeholder="用户名" type="text" autocomplete="off" />
        </el-form-item>
        <el-form-item v-if="isLogin === false" class="list-style" prop="phone">
          <el-input size="large" v-model="ruleForm.phone" placeholder="手机号" type="number" autocomplete="off" />
        </el-form-item>
        <template v-if="isLogin === true">
          <el-form-item class="list-style" prop="pass">
            <el-input size="large" v-model="ruleForm.pass" placeholder="密码" type="password" autocomplete="off" />
          </el-form-item>
        </template>
        <template v-else>
          <el-form-item class="list-style" prop="pass">
            <el-input size="large" v-model="ruleForm.pass" placeholder="密码" type="password" autocomplete="off" />
          </el-form-item>
          <el-form-item class="list-style" prop="checkPass">
            <el-input size="large" v-model="ruleForm.checkPass" placeholder="再次密码" type="password" autocomplete="off" />
          </el-form-item>
          <el-form-item class="list-style" prop="invitationCode">
            <el-input size="large" v-model="ruleForm.invitationCode" placeholder="邀请码" type="text" autocomplete="off" />
          </el-form-item>
        </template>
        <el-row :gutter="24" justify="center">
          <template v-if="isLogin === true">
            <el-col class="footer-style">
              <el-button class="login-btn" size="large" type="primary" @click="submitForm(ruleFormRef)">登录</el-button>
            </el-col>
          </template>
          <template v-else>
            <el-col class="footer-style">
              <el-button class="login-btn" size="large" type="primary" @click="submitForm(ruleFormRef)">注册</el-button>
            </el-col>
          </template>
        </el-row>
        <el-row justify="center">
          <template v-if="isLogin === true">
            <el-col class="footer-style">
              <el-button class="switch-login" type="text" @click="switchLogin()">注册账号</el-button>
            </el-col>
          </template>
          <template v-else>
            <el-col class="footer-style">
              <el-button class="switch-login" type="text" @click="switchLogin()">登录账号</el-button>
            </el-col>
          </template>
        </el-row>
      </el-form>
    </div>

<!--    <div class="company">-->
<!--      Design By NULL Group-->
<!--    </div>-->
  </div>
</template>

<script setup lang="ts">
  import { reactive, ref } from 'vue';
  import { useRouter } from 'vue-router';
  import type { FormInstance } from 'element-plus';
  import userStore from '@/store/user';
  import service from '@/api/http';
  import routerStore from "@/store/router";
  import logoUrl from '@/assets/img/logo.jpg'

  const router = useRouter();

  const isLogin = ref(true);
  const ruleFormRef = ref<FormInstance>();

  const ruleForm = reactive({
    userName: '',
    phone: '',
    pass: '',
    checkPass: '',
    invitationCode: '',
  });

  /**
   * 切换登录类型
   */
  const switchLogin = () => {
    isLogin.value = !isLogin.value;
  }

  /**
   * 验证用户名
   * @param rule
   * @param value
   * @param callback
   */
  const validateUserName = (rule: any, value: any, callback: any) => {
    if (value === '') {
      callback(new Error('请输入用户名'));
    } else {
      callback();
    }
  };

  /**
   * 验证手机号
   * @param rule
   * @param value
   * @param callback
   */
  const validatePhone = (rule: any, value: any, callback: any) => {
    if (value === '' || value.length !== 11) {
      callback(new Error('请正确输入手机号'));
    } else {
      callback();
    }
  };

  /**
   * 验证密码
   * @param rule
   * @param value
   * @param callback
   */
  const validatePass = (rule: any, value: any, callback: any) => {
    console.log(value)
    if (value === '' || value.length < 6) {
      callback(new Error('请输入6个字符以上的密码'));
    } else {
      if (ruleForm.checkPass !== '') {
        if (!ruleFormRef.value) return;
        ruleFormRef.value.validateField('checkPass', () => null);
      }
      callback();
    }
  };

  /**
   * 验证确认密码
   * @param rule
   * @param value
   * @param callback
   */
  const validatePass2 = (rule: any, value: any, callback: any) => {
    if (value === '') {
      callback(new Error('请再次输入密码'));
    } else if (value !== ruleForm.pass) {
      callback(new Error('两次密码输入不一致'));
    } else {
      callback();
    }
  };

  const validateCode = (rule: any, value: any, callback: any) => {
    if (value === '') {
      callback(new Error('请输入邀请码'));
    } else {
      callback();
    }
  };

  /**
   * 定义校验规则
   */
  const rules = reactive({
    userName: [{ validator: validateUserName, trigger: 'blur' }],
    phone: [{ validator: validatePhone, trigger: 'blur'}],
    pass: [{ validator: validatePass, trigger: 'blur' }],
    checkPass: [{ validator: validatePass2, trigger: 'blur' }],
    invitationCode: [{ validator: validateCode, trigger: 'blur' }],
  });

  /**
   * 登录注册请求回调用
   */
  const handleSign = async (submitData: object) => {
    if (isLogin.value) {
      // 注册账号
      const result = await service.post('/user/signin', submitData);
      console.log('登录：', result);
      if (result?.code === 0) {
        ElMessage({
          type: 'success',
          offset: 20,
          message: '登录成功！'
        })

        // 更新路由Store
        userStore().$patch({
          username: result?.data?.Username,
          token: result?.data?.Token,
        })
        window.localStorage.setItem("username", result?.data?.Username)
        window.localStorage.setItem("token", result?.data?.Token)

        const userInfo = await service.get('/user/info', { token: result?.data?.Token, username: result?.data?.Username });
        if (userInfo?.code === 0) {
          console.log('用户信息：', userInfo);
          userStore().$patch({
            uid: userInfo?.data?.id,
            phone: userInfo?.data?.Phone,
            inviteCode: userInfo?.data?.InviteCode,
            avatarUrl: userInfo?.data?.AvatarUrl,
            balance: userInfo?.data?.Balance,
          })

          window.localStorage.setItem("uid", userInfo?.data?.id)
          window.localStorage.setItem("phone", userInfo?.data?.Phone)
          window.localStorage.setItem("inviteCode", userInfo?.data?.InviteCode)
          window.localStorage.setItem("avatarUrl", userInfo?.data?.AvatarUrl)
          window.localStorage.setItem("balance", userInfo?.data?.Balance)
        }

        // 跳转首页
        router.push(`/home`);
      }
    } else {
      // 注册账号
      const result = await service.post('/user/signup', submitData);
      console.log(result);
      if (result?.code === 0) {
        ElMessage({
          type: 'success',
          offset: 20,
          message: '注册成功，去登录！'
        })
        isLogin.value = true;
      }
    }
  }

  /**
   * 登录/注册
   * @param formEl
   */
  const submitForm = (formEl: FormInstance | undefined) => {
    if (!formEl) return;
    // eslint-disable-next-line consistent-return
    formEl.validate((valid, invalidFields) => {
      if (valid) {
        console.log(ruleForm);

        const submitData = {
          username: ruleForm.userName,
          phone: ruleForm.phone,
          password: ruleForm.pass,
          inviteCode: isLogin.value ? undefined : ruleForm.invitationCode,
        }

        handleSign(submitData);
        console.log('提交');
      } else {
        console.log('错误提交!');
        return false;
      }
    });
  };
</script>

<style lang="scss" scoped>
  .login-bg {
    position: fixed;
    width: 100vw;
    height: 100vh;
    background-image: url('http://2022aa.external.2022AA.com/static/assets/bg.png');
    background-size: cover;
    background-repeat: no-repeat;
    background-position: 15% 100%;
    overflow: auto;
  }

  .login-body {
    display: flex;
    flex-flow: column nowrap;
    justify-content: center;
    align-items: center;
    width: 80vw;
    height: auto;
    margin: 20vh auto;
    padding: 2rem 1.5rem;
    border: 2px solid #d9d0f3;
    border-radius: 8px;
    box-shadow: 3px 3px 3px #544199;

    .top-title {
      font-size: 1.3rem;
      font-weight: bold;
      text-align: center;
      margin: 0 0 2rem;
    }

    .form-style {
      width: 100%;

      .input-style {
        height: 2rem;
      }
    }

    .list-style {
      margin-bottom: 1.5rem;
    }

    .footer-style {
      text-align: center;
    }

    .login-btn {
      font-size: 1.1rem;
      width: 12rem;
      color: #9277fd;
      margin-bottom: 0.5rem;
      background-color: #91fbcd;
      border: 2px solid #343233;
      border-radius: 8px;
    }

    .switch-login {
      color: #ddd4f3;
      text-decoration: underline;
    }
  }

  .company {
    position: fixed;
    bottom: 2.5rem;
    left: 50%;
    transform: translateX(-50%);
    font-style: italic;
  }
</style>
