<style scoped>
#loginForm {
  max-width: 420px;
  margin: 80px auto;
  padding: 24px;
  background: #fff;
  border-radius: 4px;
  box-shadow: 0 2px 12px rgba(0,0,0,.1);
}

.login-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>

<template>
  <form id="loginForm" class="layui-form" @submit.prevent="onSubmit">
    <div class="layui-form-item">
      <label class="layui-form-label">{{$L('email')}}</label>
      <div class="layui-input-block">
        <input class="layui-input" v-model="email" :placeholder="$L('email')" required>
      </div>
    </div>
    <div class="layui-form-item">
      <label class="layui-form-label">{{$L('password')}}</label>
      <div class="layui-input-block">
        <input class="layui-input" v-model="password" type="password" :placeholder="$L('password')" required>
      </div>
    </div>
    <div class="layui-form-item login-actions">
      <label>
        <input type="checkbox" v-model="remember"> {{$L('remember me')}}
      </label>
      <button class="layui-btn" type="submit">{{$L('login')}}</button>
    </div>
  </form>
</template>

<script>
export default {
  name: 'login',
  data(){
    return {
      email: '',
      password: '',
      remember: false
    }
  },
  methods: {
    onSubmit(){
      var vm = this;
      this.$rest.POST('session', {email: this.email, password: this.password, remember: this.remember}).
        onsucceed(200, (resp)=>{
          vm.$store.commit('setEmail', resp.email);
          vm.$store.commit('setRole', resp.role);
          vm.$router.push('/');
        }).
        onfailed((msg)=>{vm.$bus.$emit('error', msg)}).
        do();
    }
  }
}
</script>
