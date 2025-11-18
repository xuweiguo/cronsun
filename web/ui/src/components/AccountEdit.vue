<style scoped>
.account-edit-form {
  max-width: 640px;
}

.account-edit-form .layui-form-item {
  margin-bottom: 18px;
}

.role-options {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}
</style>

<template>
  <div>
    <div v-if="$store.getters.role === 1" class="account-edit-form">
      <form class="layui-form" @submit.prevent="submit">
        <h3>{{$L(action == 'CREATE' ? 'add account' : 'edit account')}}</h3>
        <div class="layui-form-item">
          <label class="layui-form-label">{{$L('email')}}</label>
          <div class="layui-input-block">
            <input class="layui-input" v-model="account.email" required>
          </div>
        </div>
        <div class="layui-form-item">
          <label class="layui-form-label">{{$L('password')}}</label>
          <div class="layui-input-block">
            <input class="layui-input" v-model="account.password" type="password" autocomplete="off">
          </div>
        </div>
        <div class="layui-form-item">
          <label class="layui-form-label">{{$L('role')}}</label>
          <div class="layui-input-block role-options">
            <label><input type="radio" value="1" v-model.number="account.role"> Administrator</label>
            <label><input type="radio" value="2" v-model.number="account.role"> Developer</label>
            <label><input type="radio" value="3" v-model.number="account.role"> Reporter(ReadOnly)</label>
          </div>
        </div>
        <div class="layui-form-item">
          <label class="layui-form-label">{{$L('job group')}}</label>
          <div class="layui-input-block">
            <Dropdown :title="$L('select a group')" :items="groups" :selected="account.group" @change="changeGroup" :multiple="true"/>
          </div>
        </div>
        <div class="layui-form-item" v-if="action === 'UPDATE'">
          <label class="layui-form-label">{{$L('status')}}</label>
          <div class="layui-input-block role-options">
            <label><input type="radio" value="1" v-model.number="account.status"> {{$L('active')}}</label>
            <label><input type="radio" value="-1" v-model.number="account.status"> {{$L('ban')}}</label>
          </div>
        </div>
        <div class="layui-form-item">
          <button class="layui-btn" type="submit" :disabled="loading">
            {{$L('save')}}
          </button>
        </div>
      </form>
    </div>
    <div v-else>Access deny.</div>
  </div>
</template>

<script>
import Dropdown from "./basic/Dropdown.vue";

export default {
  name: 'account-edit',
  components: {Dropdown},

  data: function(){
    return {
      loading: true,
      action: '',
      groups: [],
      account: {
        originEmail: '',
        email: '',
        password: '',
        role: 1,
        status: 1,
        group: []
      }
    }
  },

  mounted: function(){
    if (this.$route.path.indexOf('/admin/account/add') === 0) {
      this.action = 'CREATE';
      this.loading = false;
    } else {
      this.action = 'UPDATE';
      this.$rest.GET('admin/account/'+this.$route.query.email).
        onsucceed(200, (resp)=>{
          this.account = {
            originEmail: resp.email,
            email: resp.email,
            password: '',
            role: resp.role,
            status: resp.status,
            group: resp.group || []
          }
        }).
        onfailed((msg)=>{this.$bus.$emit('error', msg)}).
        onend(()=>{this.loading = false}).
        do();
    }

    this.$rest.GET('job/groups').onsucceed(200, (resp)=>{
      !resp.includes('default') && resp.unshift('default');
      this.groups = resp;
    }).do();
  },

  methods: {
    submit(){
      if (this.loading) return;
      var req, code;
      this.loading = true;
      if (this.action === 'CREATE') {
        req = this.$rest.PUT('admin/account', this.account);
        code =  204;
      } else {
        req = this.$rest.POST('admin/account', this.account)
        code =  200;
      }

      req.onsucceed(code, ()=>{
          this.$router.push('/admin/account/list');
        }).
        onfailed((msg)=>{this.$bus.$emit('error', msg)}).
        onend(()=>{this.loading = false}).
        do();
    },
    changeGroup: function(val){
      this.account.group = val;
    },
  }
}
</script>
