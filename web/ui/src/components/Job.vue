<style scoped>
.job-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  flex-wrap: wrap;
  gap: 12px;
}

.job-toolbar__actions {
  display: flex;
  gap: 8px;
  align-items: center;
  flex-wrap: wrap;
}

.job-filter-card {
  margin-bottom: 16px;
}

.job-filter-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 16px;
}

.job-form-item {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.job-table-wrapper {
  position: relative;
}

.job-loading {
  text-align: center;
  padding: 24px 0;
  color: #999;
}

.job-name-link {
  color: #1e9fff;
}

.job-operations {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}
</style>

<template>
  <div class="job-page">
    <div class="job-toolbar">
      <div class="job-toolbar__actions">
        <button class="layui-btn layui-btn-primary" type="button" @click="$router.push('/job/executing')">
          {{$L('view executing jobs')}}
        </button>
        <button class="layui-btn layui-btn-primary layui-btn-sm" type="button" @click="refresh">{{$L('refresh')}}</button>
        <div class="layui-btn-group">
          <button class="layui-btn layui-btn-sm" type="button" @click="batched = !batched">
            {{$L('batch')}}
          </button>
          <button
            v-if="batched"
            class="layui-btn layui-btn-sm"
            type="button"
            :disabled="batchIds.length === 0"
            @click="batch('start')">
            {{$L('open')}}
          </button>
          <button
            v-if="batched"
            class="layui-btn layui-btn-sm layui-btn-warm"
            type="button"
            :disabled="batchIds.length === 0"
            @click="batch('pause')">
            {{$L('pause')}}
          </button>
        </div>
      </div>
      <button class="layui-btn" type="button" @click="$router.push('/job/create')">
        {{$L('create job')}}
      </button>
    </div>

    <div class="layui-card job-filter-card">
      <div class="layui-card-body job-filter-grid">
        <div class="job-form-item">
          <label>{{$L('group filter')}}</label>
          <Dropdown :title="$L('select a group')" :items="groups" @change="changeGroup" :selected="group"/>
        </div>
        <div class="job-form-item">
          <label>{{$L('node filter')}}</label>
          <Dropdown :title="$L('select a node')" :items="nodes" @change="changeNode" :selected="node"/>
        </div>
      </div>
    </div>

    <div class="job-table-wrapper">
      <table class="layui-table" v-if="!loading && jobs.length">
        <thead>
          <tr>
            <th>{{$L('operation')}}</th>
            <th>{{$L('status')}}</th>
            <th>{{$L('group')}}</th>
            <th>{{$L('user')}}</th>
            <th>{{$L('name')}}</th>
            <th>{{$L('latest executed')}}</th>
            <th>{{$L('executing result')}}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(row, index) in jobs" :key="row.group + '-' + row.id">
            <td>
              <div v-if="batched">
                <input type="checkbox" :value="row.group + '/' + row.id" v-model="batchIds">
              </div>
              <div v-else class="job-operations">
                <button class="layui-btn layui-btn-xs" @click="editJob(row.group, row.id)">{{$L('edit')}}</button>
                <button class="layui-btn layui-btn-xs" @click="changeStatus(row.group, row.id, index, !row.pause)">
                  {{row.pause ? $L('open') : $L('pause')}}
                </button>
                <button class="layui-btn layui-btn-danger layui-btn-xs" @click="removeJob(row.group, row.id, index)">{{$L('delete')}}</button>
              </div>
            </td>
            <td>{{row.pause ? $L('pause') : $L('open')}}</td>
            <td>{{row.group}}</td>
            <td>{{row.user && row.user.length > 0 ? row.user : '-'}}</td>
            <td>
              <router-link class="job-name-link" :to="buildEditJobURL(row.group, row.id)">{{row.name}}</router-link>
            </td>
            <td>
              <div v-if="!row.latestStatus">-</div>
              <div v-else>
                {{formatLatest(row.latestStatus)}}
                <div><small>{{formatNextRunTime(row.nextRunTime)}}</small></div>
              </div>
            </td>
            <td>
              <span v-if="!row.latestStatus">-</span>
              <span v-else>
                <router-link :to="'/log/'+row.latestStatus.refLogId">
                  {{$L(row.latestStatus.success ? 'successed' : 'failed')}}
                </router-link>
                |
                <router-link :to="{path: 'log', query: {latest:true, ids: row.id}}">latest</router-link>
                |
                <button class="layui-btn layui-btn-primary layui-btn-xs" type="button" @click.prevent="showExecuteJobModal(row.name, row.group, row.id)">
                  {{$L('execute')}}
                </button>
              </span>
            </td>
          </tr>
        </tbody>
      </table>
      <div class="job-loading" v-else>
        {{loading ? $L('loading') : $L('no data')}}
      </div>
    </div>
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
