<template>
  <el-form v-else label-position="top" :model="account" v-loading="loading" class="account-edit-form">
    <h3>{{$L(action == 'CREATE' ? 'add account' : 'edit account')}}</h3>
    <el-form-item :label="$L('email')">
      <el-input v-model="account.email"/>
    </el-form-item>
    <el-form-item :label="$L('password')">
      <el-input v-model="account.password" type="password" autocomplete="off"/>
    </el-form-item>
    <el-form-item :label="$L('role')">
      <el-radio-group v-model="account.role">
        <el-radio-button :label="1">Administrator</el-radio-button>
        <el-radio-button :label="2">Developer</el-radio-button>
        <el-radio-button :label="3">Reporter(ReadOnly)</el-radio-button>
      </el-radio-group>
    </el-form-item>
    <el-form-item :label="$L('job group')">
      <Dropdown :title="$L('select a group')" :items="groups" :selected="account.group" @change="changeGroup" :multiple="true"/>
    </el-form-item>
    <el-form-item v-if="action === 'UPDATE'" :label="$L('status')">
      <el-radio-group v-model="account.status">
        <el-radio-button :label="1">{{$L('active')}}</el-radio-button>
        <el-radio-button :label="-1">{{$L('ban')}}</el-radio-button>
      </el-radio-group>
    </el-form-item>
    <el-form-item>
      <el-button type="primary" @click="submit" :loading="loading">{{$L('save')}}</el-button>
    </el-form-item>
  </el-form>
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
