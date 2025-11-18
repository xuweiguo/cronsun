<style scoped>
.account-toolbar {
  margin-bottom: 16px;
  text-align: right;
}

.account-table {
  width: 100%;
}
</style>

<template>
  <div>
    <div v-if="$store.getters.role === 1">
      <div class="account-toolbar">
        <button class="layui-btn" type="button" @click="$router.push('/admin/account/add')">
          {{$L('add account')}}
        </button>
      </div>
      <table class="layui-table account-table">
        <thead>
          <tr>
            <th>{{$L('email')}}</th>
            <th>{{$L('role')}}</th>
            <th>{{$L('status')}}</th>
            <th>{{$L('added time')}}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="row in list" :key="row.email">
            <td>
              <router-link :to="'/admin/account/edit?email='+row.email">{{row.email}}</router-link>
            </td>
            <td>{{roleMap(row.role)}}</td>
            <td>{{$L(statusMap(row.status))}}</td>
            <td>{{row.createTime}}</td>
          </tr>
          <tr v-if="!list || list.length === 0">
            <td colspan="4">{{$L('no data')}}</td>
          </tr>
        </tbody>
      </table>
    </div>
    <div v-else>Access deny.</div>
  </div>
</template>

<script>
export default {
  name: 'account',

  data(){
    return {
      list: []
    }
  },

  mounted(){
    this.$rest.GET('/admin/accounts').
      onsucceed(200, (resp)=>{
        this.list = resp;
      }).
      do();
  },

  methods: {
    statusMap(s) {
      switch (s){
        case -1: return 'banned';
        case 1: return 'actived';
      }
      return '';
    },

    roleMap(s) {
      switch (s){
        case 1: return 'Administrator';
        case 2: return 'Developer';
        case 3: return 'Repoter';
      }
      return '';
    }
  }
}
</script>
