<style scoped>
@media (min-width: 600px) {
  #loginForm {
    width: 500px;
    margin: 100px auto 0;
  }
}
</style>

<template>
  <div>
    <el-form id="loginForm" label-position="top" @submit.native.prevent="onSubmit">
      <el-form-item :label="$L('email')">
        <el-input v-model="email" :placeholder="$L('email')"/>
      </el-form-item>
      <el-form-item :label="$L('password')">
        <el-input v-model="password" type="password" :placeholder="$L('password')"/>
      </el-form-item>
      <el-form-item>
        <el-checkbox v-model="remember">{{$L('remember me')}}</el-checkbox>
      </el-form-item>
      <el-button type="primary" @click="onSubmit">{{$L('login')}}</el-button>
    </el-form>
  </div>
</template>

<script>
export default {
  name: 'login',

  data: function(){
    return {
      email: '',
      password: '',
      remember: false
    }
  },

  methods: {
    onSubmit(){
      var vm = this;
      this.$rest.GET('session?email='+this.email+'&password='+this.password+'&remember='+this.remember).
      onsucceed(200, (resp)=>{
        vm.$store.commit('setEmail', resp.email);
        vm.$store.commit('setRole', resp.role);
        vm.$store.commit('enabledAuth', resp.enabledAuth);
        vm.$router.push('/')
        vm.$loadConfiguration();
      }).
      do();
    }
  }
}
</script>
