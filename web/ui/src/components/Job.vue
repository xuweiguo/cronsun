<style scoped>
.job-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.job-toolbar__actions {
  display: flex;
  gap: 8px;
  align-items: center;
}

.job-filter-card {
  margin-bottom: 16px;
}

.job-table {
  margin-top: 8px;
}

.job-name-link {
  color: #409EFF;
}

.job-name-link:hover {
  text-decoration: underline;
}
</style>

<template>
  <div class="job-page">
    <div class="job-toolbar">
      <div class="job-toolbar__actions">
        <el-button type="primary" icon="el-icon-view" @click="$router.push('/job/executing')">
          {{$L('view executing jobs')}}
        </el-button>
        <el-button icon="el-icon-refresh" circle @click="refresh"></el-button>
        <el-button-group>
          <el-button icon="el-icon-collection" @click="batched = !batched">
            {{$L('batch')}}
          </el-button>
          <el-button
            v-if="batched"
            type="success"
            icon="el-icon-video-play"
            :disabled="batchIds.length === 0"
            @click="batch('start')"/>
          <el-button
            v-if="batched"
            type="warning"
            icon="el-icon-video-pause"
            :disabled="batchIds.length === 0"
            @click="batch('pause')"/>
        </el-button-group>
      </div>
      <el-button type="primary" icon="el-icon-plus" @click="$router.push('/job/create')">
        {{$L('create job')}}
      </el-button>
    </div>

    <el-card class="job-filter-card" shadow="never">
      <el-row :gutter="16">
        <el-col :sm="12" :xs="24">
          <label>{{$L('group filter')}}</label>
          <Dropdown :title="$L('select a group')" :items="groups" @change="changeGroup" :selected="group"/>
        </el-col>
        <el-col :sm="12" :xs="24">
          <label>{{$L('node filter')}}</label>
          <Dropdown :title="$L('select a node')" :items="nodes" @change="changeNode" :selected="node"/>
        </el-col>
      </el-row>
    </el-card>

    <el-table
      class="job-table"
      :data="jobs"
      stripe
      border
      v-loading="loading"
      :empty-text="$L('no data')">
      <el-table-column
        :label="$L('operation')"
        width="140"
        fixed="left">
        <template slot-scope="{row, $index}">
          <el-checkbox
            v-if="batched"
            v-model="batchIds"
            :label="row.group + '/' + row.id">
          </el-checkbox>
          <el-dropdown v-else trigger="click" @command="(cmd) => handleAction(row, $index, cmd)">
            <span class="el-dropdown-link">
              <i class="el-icon-more"></i>
            </span>
            <el-dropdown-menu slot="dropdown">
              <el-dropdown-item command="edit">{{$L('edit')}}</el-dropdown-item>
              <el-dropdown-item command="toggle">
                {{row.pause ? $L('open') : $L('pause')}}
              </el-dropdown-item>
              <el-dropdown-item divided command="delete" style="color:#F56C6C;">
                {{$L('delete')}}
              </el-dropdown-item>
            </el-dropdown-menu>
          </el-dropdown>
        </template>
      </el-table-column>
      <el-table-column :label="$L('status')" width="120">
        <template slot-scope="{row}">
          <el-tag :type="row.pause ? 'warning' : 'success'">
            {{row.pause ? $L('pause') : $L('open')}}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="group" :label="$L('group')" width="140"></el-table-column>
      <el-table-column :label="$L('user')" width="160">
        <template slot-scope="{row}">
          {{row.user && row.user.length > 0 ? row.user : '-'}}
        </template>
      </el-table-column>
      <el-table-column :label="$L('name')" min-width="200">
        <template slot-scope="{row}">
          <router-link class="job-name-link" :to="buildEditJobURL(row.group, row.id)">{{row.name}}</router-link>
        </template>
      </el-table-column>
      <el-table-column :label="$L('latest executed')" min-width="220">
        <template slot-scope="{row}">
          <div v-if="!row.latestStatus">-</div>
          <div v-else>{{formatLatest(row.latestStatus)}}</div>
          <small>{{formatNextRunTime(row.nextRunTime)}}</small>
        </template>
      </el-table-column>
      <el-table-column :label="$L('executing result')" min-width="280">
        <template slot-scope="{row}">
          <span v-if="!row.latestStatus">-</span>
          <span v-else>
            <router-link :to="'/log/'+row.latestStatus.refLogId">
              {{$L(row.latestStatus.success ? 'successed' : 'failed')}}
            </router-link>
            |
            <router-link :to="{path: 'log', query: {latest:true, ids: row.id}}">latest</router-link>
            |
            <el-tooltip :content="$L('click to select a node and re-execute job')">
              <el-button
                type="text"
                icon="el-icon-refresh"
                @click.prevent="showExecuteJobModal(row.name, row.group, row.id)">
              </el-button>
            </el-tooltip>
          </span>
        </template>
      </el-table-column>
    </el-table>
    <ExecuteJob ref="executeJobModal"/>
  </div>
</template>

<script>
import Dropdown from './basic/Dropdown.vue';
import ExecuteJob from './ExecuteJob.vue';
import {formatTime, formatDuration} from '../libraries/functions';

export default {
  name: 'job',
  data: function(){
    return {
      batched: false,
      batchIds: [],
      groups: [],
      group: '',
      nodes: [],
      node: '',
      jobs: [],
      loading: false
    }
  },

  mounted: function(){
    this.fillParams();
    var vm = this;

    this.$rest.GET('job/groups').onsucceed(200, (resp)=>{
      !resp.includes('default') && resp.unshift('default');
      resp.unshift({value: '', name: vm.$L('all groups')});
      vm.groups = resp;
      this.fetchList(this.buildQuery());
    }).do();

    var nodes = Array.from(this.$store.getters.dropdownNodes);
    nodes.unshift({value: '', name: this.$L('all nodes')});
    vm.nodes = nodes;
  },

  watch: {
    '$route': function(){
      this.fillParams();
      this.fetchList(this.buildQuery());
    }
  },

  methods: {
    fillParams: function(){
      this.group = this.$route.query.group || '';
      this.node = this.$route.query.node || '';
    },

    changeGroup: function(val){
      this.group = val;
      this.$router.push('job?'+this.buildQuery());
    },

    changeNode: function(val){
      this.node = val;
      this.$router.push('job?'+this.buildQuery());
    },

    buildQuery: function(){
      var params = [];
      if (this.group) params.push('group='+this.group);
      if (this.node) params.push('node='+this.node);
      return params.join('&');
    },

    fetchList: function(query){
      var vm = this;
      this.loading = true;
      this.$rest.GET('jobs?'+query).onsucceed(200, (resp)=>{
        vm.jobs = resp;
      }).onend(()=>{
        vm.loading = false;
      }).do();
    },

    refresh: function(){
      this.fetchList(this.buildQuery());
    },

    editJob: function(group, id){
        this.$router.push( this.buildEditJobURL(group, id) );
    },

    buildEditJobURL: function(group, id){
        var query = this.buildQuery();
        if(query == ''){
            return '/job/edit/'+group+'/'+id;
        }else{
            return '/job/edit/'+group+'/'+id+ '?'+query;
        }
    },

    removeJob: function(group, id, index){
      var vm = this;
      this.$rest.DELETE('job/'+group+'-'+id).onsucceed(204, (resp)=>{
        vm.jobs.splice(index, 1);
      }).do();
    },

    changeStatus: function(group, id, index, isPause){
      var vm = this;
      this.$rest.POST('job/'+group+'-'+id, {"pause": isPause}).onsucceed(200, (resp)=>{
        vm.refresh();
      }).do();
    },

    handleAction(row, index, command){
      if (command === 'edit') {
        this.editJob(row.group, row.id);
      } else if (command === 'toggle') {
        this.changeStatus(row.group, row.id, index, !row.pause);
      } else if (command === 'delete') {
        this.removeJob(row.group, row.id, index);
      }
    },

    formatExecResult: function(st){
      if (!st) return '-';
      return
    },

    formatLatest: function(latest){
      return this.$L('on {node} took {times}, {begin ~ end}', this.$store.getters.hostshowsWithoutTip(latest.node), formatDuration(latest.beginTime, latest.endTime), formatTime(latest.beginTime, latest.endTime));
    },

    formatNextRunTime: function(nextRunTime){
      return this.$L('next schedule: {nextTime}', nextRunTime);
    },

    showExecuteJobModal: function(jobName, jobGroup, jobId){
      this.$refs.executeJobModal.show(jobName, jobGroup, jobId);
    },

    batch: function(op){
      switch(op) {
        case 'start': break;
        case 'pause': break;
        default: return;
      }

      var vm = this;
      this.$rest.POST('jobs/'+op, this.batchIds).onsucceed(200, (resp)=>{
        vm.refresh();
        vm.batchIds = [];
      }).onfailed((msg)=>{
        vm.$bus.$emit('warning', msg);
      }).do();
    }
  },

  components: {
    Dropdown,
    ExecuteJob
  }
}
</script>
