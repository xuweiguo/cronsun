<style scoped>
.job-executing-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}
</style>

<template>
  <div>
    <div class="job-executing-toolbar">
      <div>
        <el-button type="primary" icon="el-icon-s-operation" @click="$router.push('/job')">
          {{$L('view job list')}}
        </el-button>
        <el-button icon="el-icon-refresh" circle @click="refresh"></el-button>
      </div>
      <el-button type="primary" icon="el-icon-plus" @click="$router.push('/job/create')">
        {{$L('create job')}}
      </el-button>
    </div>
    <el-card shadow="never" :body-style="{padding: '16px'}">
      <el-form label-position="top" v-loading="loading">
        <el-form-item :label="$L('job ID')">
          <el-input type="text" v-model="ids" :placeholder="$L('multiple IDs can separated by commas')"/>
        </el-form-item>
        <el-form-item :label="$L('select groups')">
          <Dropdown :title="$L('select groups')" :items="prefetchs.groups" @change="changeGroup" :selected="groups" :multiple="true"/>
        </el-form-item>
        <el-form-item :label="$L('select nodes')">
          <Dropdown :title="$L('select nodes')" :items="$store.getters.dropdownNodes" @change="changeNodes" :selected="nodes" :multiple="true"/>
        </el-form-item>
        <div style="text-align:right;">
          <el-button type="primary" @click="submit">{{$L('submit query')}} </el-button>
        </div>
      </el-form>
    </el-card>
    <el-table :data="executings" border stripe style="margin-top:16px;" :empty-text="$L('no data')">
      <el-table-column :label="$L('job name')" min-width="200">
        <template slot-scope="{row}">
          <router-link :to="'/job/edit/'+row.group+'/'+row.jobId">{{row.jobName}}</router-link>
        </template>
      </el-table-column>
      <el-table-column prop="group" :label="$L('job group')" width="180"></el-table-column>
      <el-table-column :label="$L('node')" min-width="180">
        <template slot-scope="{row}">{{$store.getters.hostshows(row.nodeId)}}</template>
      </el-table-column>
      <el-table-column prop="id" :label="$L('process ID')" width="160"></el-table-column>
      <el-table-column prop="time" :label="$L('starting time')" width="200"></el-table-column>
      <el-table-column :label="$L('operation')" width="160">
        <template slot-scope="{row, $index}">
          <el-button type="danger" plain size="mini" @click="killProc(row, $index)">{{$L('kill process')}}</el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script>
import Dropdown from './basic/Dropdown.vue';
import {split} from '../libraries/functions';

export default {
  name: 'job-executing',
  data(){
    return {
      prefetchs: {groups: []},
      loading: false,
      groups: [],
      ids: '',
      nodes: [],
      executings: []
    }
  },

  mounted(){
    this.groups = split(this.$route.query.groups, ',');
    this.nodes = split(this.$route.query.nodes, ',');
    this.ids = this.$route.query.ids || '';

    this.$rest.GET('job/groups').onsucceed(200, (resp)=>{
      !resp.includes('default') && resp.unshift('default');
      this.prefetchs.groups = resp;
      this.fetchList(this.buildQuery());
    }).do();
  },

  watch: {
    '$route': function(){
      this.groups = split(this.$route.query.groups, ',');
      this.nodes = split(this.$route.query.nodes, ',');
      this.ids = this.$route.query.ids || '';
      this.fetchList(this.buildQuery());
    }
  },

  methods: {
    changeGroup(val){
      this.groups = Array.isArray(val) ? val : split(val, ',');
    },

    changeNodes(val){
      this.nodes = Array.isArray(val) ? val : split(val, ',');
    },

    submit(){
      this.$router.push('/job/executing?'+this.buildQuery());
    },

    killProc(proc, index) {
      this.$confirm(this.$L('whether to kill the process'), this.$L('operation'), {
        confirmButtonText: this.$L('confirm'),
        cancelButtonText: this.$L('cancel'),
        type: 'warning'
      }).then(() => {
        var params = []
        params.push('node='+proc.nodeId)
        params.push('group='+proc.group)
        params.push('job='+proc.jobId)
        params.push('pid='+proc.id)
        this.$rest.DELETE('job/executing?' + params.join('&'))
          .onsucceed(200, () => this.$bus.$emit('success', this.$L('command has been sent to the node')))
          .onfailed((resp) => this.$bus.$emit('error', resp))
          .do();
      }).catch(()=>{});
    },

    buildQuery(){
      var params = [];
      if (this.groups && this.groups.length > 0) params.push('groups='+this.groups.join(','));
      if (this.nodes && this.nodes.length > 0) params.push('nodes='+this.nodes.join(','));
      if (this.ids) params.push('ids='+this.ids);
      return params.join('&');
    },

    fetchList(query){
      this.loading = true;
      this.$rest.GET('job/executing?'+query).
      onsucceed(200, (resp)=>{
        this.executings = resp;
      }).
      onend(()=>{this.loading = false}).
      do();
    },

    refresh(){
      this.fetchList(this.buildQuery());
    }
  },

  components: {
    Dropdown
  }
}
</script>
