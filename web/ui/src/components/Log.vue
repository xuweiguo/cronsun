<style scoped>
.log-filter-card {
  margin-bottom: 16px;
  background: #fff;
  border-radius: 4px;
  padding: 16px;
}

.log-filter__actions {
  text-align: right;
  margin-top: 12px;
}

.log-filter-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 16px;
}

.log-warning {
  color: #FF5722;
  font-weight: 600;
}
</style>

<template>
  <div class="log-page">
    <div class="log-filter-card">
      <form class="layui-form" @submit.prevent="submit" @keyup.enter="submit">
        <div class="log-filter-grid">
          <div>
            <label>{{$L('job name')}}</label>
            <input class="layui-input" v-model="names" :placeholder="$L('multiple names can separated by commas')"/>
          </div>
          <div>
            <label>{{$L('job ID')}}</label>
            <input class="layui-input" v-model="ids" :placeholder="$L('multiple IDs can separated by commas')"/>
          </div>
          <div>
            <label>{{$L('node')}}</label>
            <input
              class="layui-input"
              v-if="$store.getters.showWithHostname"
              v-model="hostnames"
              :placeholder="$L('multiple Hostnames can separated by commas')"/>
            <input
              class="layui-input"
              v-else
              v-model="ips"
              :placeholder="$L('multiple IPs can separated by commas')"/>
          </div>
          <div>
            <label>{{$L('starting date')}}</label>
            <input type="date" class="layui-input" v-model="begin"/>
          </div>
          <div>
            <label>{{$L('end date')}}</label>
            <input type="date" class="layui-input" v-model="end"/>
          </div>
        </div>
        <div class="log-filter-grid" style="margin-top:16px;">
          <label>
            <input type="checkbox" v-model="latest"> {{$L('latest result of each job on each node')}}
          </label>
          <label>
            <input type="checkbox" v-model="failedOnly"> {{$L('failure only')}}
          </label>
        </div>
        <div class="log-filter__actions">
          <button class="layui-btn" type="submit" :disabled="loading">{{$L('submit query')}}</button>
        </div>
      </form>
    </div>

    <table class="layui-table" v-if="list && list.length">
      <thead>
        <tr>
          <th>{{$L('job name')}}</th>
          <th>{{$L('executing node')}}</th>
          <th>{{$L('executing user')}}</th>
          <th>{{$L('executing time')}}</th>
          <th>{{$L('executing result')}}</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="row in list" :key="row.id">
          <td>
            <router-link :to="'/job/edit/'+row.jobGroup+'/'+row.jobId">{{row.name}}</router-link>
          </td>
          <td>
            <span :title="row.node">{{$store.getters.hostshows(row.node)}}</span>
          </td>
          <td>{{row.user}}</td>
          <td>
            <span :class="{'log-warning': durationAttention(row.beginTime, row.endTime)}">
              {{formatTime(row)}}
            </span>
          </td>
          <td>
            <router-link :to="'/log/'+row.id">
              {{$L(row.success ? 'successed' : 'failed')}}
            </router-link>
            |
            <button class="layui-btn layui-btn-xs" type="button" :title="$L('click to select a node and re-execute job')" @click.prevent="showExecuteJobModal(row.name, row.jobGroup, row.jobId)">{{$L('execute')}}</button>
          </td>
        </tr>
      </tbody>
    </table>
    <div v-else class="log-warning">{{loading ? $L('loading') : $L('no data')}}</div>
    <Pager v-if="list && list.length>0" :total="total" :maxBtn="5"/>
    <ExecuteJob ref="executeJobModal"/>
  </div>
</template>

<script>
import Pager from './basic/Pager.vue';
import ExecuteJob from './ExecuteJob.vue';
import {formatTime, formatDuration} from '../libraries/functions';

export default {
  name: 'log',
  data: function(){
    return {
      loading: false,
      names: '',
      ids: '',
      hostnames: '',
      ips: '',
      begin: '',
      end: '',
      latest: false,
      failedOnly: false,
      list: [],
      total: 0,
      page: 1
    }
  },

  mounted: function(){
      this.fillParams();
      this.fetchList(this.buildQuery());

      const time = new Date();
      const day = ("0" + time.getDate()).slice(-2);
      const month = ("0" + (time.getMonth() + 1)).slice(-2);
      const today = time.getFullYear() + "-" + (month) + "-" + (day);
      if (!this.begin && !this.end) {
        this.begin = this.end = today;
      }
  },

  watch: {
    '$route': function(){
      this.fillParams();
      this.fetchList(this.buildQuery());
    }
  },

  methods: {
    fillParams(){
      this.names = this.$route.query.names || '';
      this.ids = this.$route.query.ids || '';
      this.hostnames = this.$route.query.hostnames || '';
      this.ips = this.$route.query.ips || '';
      this.begin = this.$route.query.begin || this.begin;
      this.end = this.$route.query.end || this.end;
      this.page = this.$route.query.page || 1;
      this.latest = this.$route.query.latest === 'true' || this.$route.query.latest === true;
      this.failedOnly = this.$route.query.failedOnly === 'true' || this.$route.query.failedOnly === true;
    },

    fetchList(query){
      this.loading = true;
      var vm = this;
      this.$rest.GET('logs?'+query)
        .onsucceed(200, (resp)=>{
          vm.list = resp.list;
          vm.total = resp.total;
        })
        .onfailed((msg)=>{vm.$bus.$emit('error', msg)})
        .onend(()=>{vm.loading=false})
        .do();
    },

    buildQuery(){
      var params = [];
      if (this.names) params.push('names='+this.names);
      if (!this.$store.getters.showWithHostname && this.ids) params.push('ids='+this.ids);
      if (this.$store.getters.showWithHostname && this.hostnames) params.push('hostnames='+this.hostnames);
      if (this.ips) params.push('ips='+this.ips);
      if (this.begin) params.push('begin='+this.begin);
      if (this.end) params.push('end='+this.end);
      if (this.failedOnly) params.push('failedOnly=true');
      if (this.page == 0) this.page = 1;
      params.push('page='+this.page);
      if (this.latest) params.push('latest=true');
      return params.join('&');
    },

    submit: function(){
      var query = this.buildQuery()
      var url = '/log?'+query;
      if (this.$route.fullPath == url) {
        this.fetchList(query);
        return;
      }
      this.$router.push(url);
    },

    durationAttention: function(beginTime, endTime){
      var d = new Date(endTime) - new Date(beginTime);
      return d > 3600000*6;
    },

    formatTime: function(log){
      return this.$L('took {times}, {begin ~ end}', formatDuration(log.beginTime, log.endTime), formatTime(log.beginTime, log.endTime));
    },

    showExecuteJobModal: function(jobName, jobGroup, jobId){
      this.$refs.executeJobModal.show(jobName, jobGroup, jobId);
    }
  },
  components: {
    Pager,
    ExecuteJob
  }
}
</script>
