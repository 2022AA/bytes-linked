import { defineStore } from 'pinia';

const userStore = defineStore({
    // 这里的id必须为唯一ID
    id: 'user',
    state: () => {
        return {
            uid: '',
            username: '',
            phone: '',
            password: '',
            token: '',
            inviteCode: '',
            avatarUrl: '',
            balance: '',
        };
    },
    // 等同于vuex的getter
    getters: {
        getUserUid: (state) => state.uid,
        getUserName: (state) => state.username,
        getPhone: (state) => state.phone,
        getPassword: (state) => state.password,
        getToken: (state) => state.token,
        getInviteCode: (state) => state.inviteCode,
        getAvatarUrl: (state) => state.avatarUrl,
        getBalance: (state) => state.balance,
    },
    // pinia 放弃了 mutations 只使用 actions
    actions: {
        // actions可以用async做成异步形式
        setUserName(type: string) {
            this.username = type;
        },
        setPassword(type: string) {
            this.password = type;
        },
        setToken(type: string) {
            this.token = type;
        },
    },
});

export default userStore;
