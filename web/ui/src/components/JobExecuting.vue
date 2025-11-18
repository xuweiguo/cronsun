<style scoped>
.job-executing-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  flex-wrap: wrap;
  gap: 12px;
}

.job-form {
  background: #fff;
  padding: 16px;
  border-radius: 4px;
}

.job-form .layui-form-item {
  margin-bottom: 16px;
}

.job-table {
  margin-top: 16px;
}

.job-table td button {
  margin-right: 4px;
}
</style>

<template>
  <div>
    <div class="job-executing-toolbar">
      <div class="layui-btn-group">
        <button class="layui-btn layui-btn-primary" type="button" @click="$router.push('/job')">
          {{$L('view job list')}}
        </button>
        <button class="layui-btn layui-btn-sm" type="button" @click="refresh">{{$L('refresh')}}</button>
      </div>
      <button class="layui-btn" type="button" @click="$router.push('/job/create')">
        {{$L('create job')}}
      </button>
    </div>
    <div class="job-form">
      <form class="layui-form" @submit.prevent="submit">
        <div class="layui-form-item">
          <label class="layui-form-label">{{$L('job ID')}}</label>
          <div class="layui-input-block">
            <input type="text" class="layui-input" v-model="ids" :placeholder="$L('multiple IDs can separated by commas')">
          </div>
        </div>
        <div class="layui-form-item">
          <label class="layui-form-label">{{$L('select groups')}}</label>
          <div class="layui-input-block">
            <Dropdown :title="$L('select groups')" :items="prefetchs.groups" @change="changeGroup" :selected="groups" :multiple="true"/>
          </div>
        </div>
        <div class="layui-form-item">
          <label class="layui-form-label">{{$L('select nodes')}}</label>
          <div class="layui-input-block">
            <Dropdown :title="$L('select nodes')" :items="$store.getters.dropdownNodes" @change="changeNodes" :selected="nodes" :multiple="true"/>
          </div>
        </div>
        <div class="layui-form-item" style="text-align:right;">
          <button class="layui-btn" type="submit" :disabled="loading">{{$L('submit query')}}</button>
        </div>
      </form>
    </div>
    <table class="layui-table job-table">
      <thead>
        <tr>
          <th>{{$L('job name')}}</th>
          <th>{{$L('job group')}}</th>
          <th>{{$L('node')}}</th>
          <th>{{$L('process ID')}}</th>
          <th>{{$L('starting time')}}</th>
          <th>{{$L('operation')}}</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(row, index) in executings" :key="row.id + row.nodeId">
          <td>
            <router-link :to="'/job/edit/'+row.group+'/'+row.jobId">{{row.jobName}}</router-link>
          </td>
          <td>{{row.group}}</td>
          <td>{{$store.getters.hostshows(row.nodeId)}}</td>
          <td>{{row.id}}</td>
          <td>{{row.time}}</td>
          <td>
            <button class="layui-btn layui-btn-danger layui-btn-xs" type="button" @click="killProc(row, index)">{{$L('kill process')}}</button>
          </td>
        </tr>
        <tr v-if="executings.length === 0">
          <td colspan="6">{{$L(loading ? 'loading' : 'no data')}}</td>
        </tr>
      </tbody>
    </table>
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
      const layer = this.$layer || (this.$layui && this.$layui.layer);
      const confirmAction = () => {
        var params = []
        params.push('node='+proc.nodeId)
        params.push('group='+proc.group)
        params.push('job='+proc.jobId)
        params.push('pid='+proc.id)
        this.$rest.DELETE('job/executing?' + params.join('&'))
          .onsucceed(200, () => this.$bus.$emit('success', this.$L('command has been sent to the node')))
          .onfailed((resp) => this.$bus.$emit('error', resp))
          .do();
      };

      if (layer && layer.confirm) {
        layer.confirm(this.$L('whether to kill the process'), {title: this.$L('operation')}, (indexLayer) => {
          layer.close(indexLayer);
          confirmAction();
        });
        return;
      }
      if (window.confirm(this.$L('whether to kill the process'))) {
        confirmAction();
      }
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
