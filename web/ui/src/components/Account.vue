<style scoped>
.account-toolbar {
  margin-bottom: 16px;
  text-align: right;
}
</style>

<template>
  <div>
    <div v-if="$store.getters.role === 1">
      <div class="account-toolbar">
        <el-button type="primary" icon="el-icon-user-solid" @click="$router.push('/admin/account/add')">
          {{$L('add account')}}
        </el-button>
      </div>
      <el-table :data="list" border stripe :empty-text="$L('no data')">
        <el-table-column :label="$L('email')" min-width="220">
          <template slot-scope="{row}">
            <router-link :to="'/admin/account/edit?email='+row.email">{{row.email}}</router-link>
          </template>
        </el-table-column>
        <el-table-column :label="$L('role')" width="180">
          <template slot-scope="{row}">{{roleMap(row.role)}}</template>
        </el-table-column>
        <el-table-column :label="$L('status')" width="160">
          <template slot-scope="{row}">{{$L(statusMap(row.status))}}</template>
        </el-table-column>
        <el-table-column prop="createTime" :label="$L('added time')" width="200"></el-table-column>
      </el-table>
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
