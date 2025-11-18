<style scoped>
.log-filter-card {
  margin-bottom: 16px;
}

.log-filter__actions {
  display: flex;
  justify-content: flex-end;
}

.log-warning {
  color: #E6A23C;
  font-weight: 600;
}
</style>

<template>
  <div class="log-page">
    <el-card class="log-filter-card" shadow="never">
      <el-form label-position="top" @keyup.enter.native="submit">
        <el-row :gutter="16">
          <el-col :sm="12" :xs="24">
            <el-form-item :label="$L('job name')">
              <el-input v-model="names" :placeholder="$L('multiple names can separated by commas')"/>
            </el-form-item>
          </el-col>
          <el-col :sm="12" :xs="24">
            <el-form-item :label="$L('job ID')">
              <el-input v-model="ids" :placeholder="$L('multiple IDs can separated by commas')"/>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item :label="$L('node')">
          <el-input
            v-if="$store.getters.showWithHostname"
            v-model="hostnames"
            :placeholder="$L('multiple Hostnames can separated by commas')"/>
          <el-input
            v-else
            v-model="ips"
            :placeholder="$L('multiple IPs can separated by commas')"/>
        </el-form-item>
        <el-row :gutter="16">
          <el-col :sm="12" :xs="24">
            <el-form-item :label="$L('starting date')">
              <el-date-picker type="date" value-format="yyyy-MM-dd" v-model="begin" style="width: 100%"/>
            </el-form-item>
          </el-col>
          <el-col :sm="12" :xs="24">
            <el-form-item :label="$L('end date')">
              <el-date-picker type="date" value-format="yyyy-MM-dd" v-model="end" style="width: 100%"/>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="16">
          <el-col :sm="12" :xs="24">
            <el-checkbox v-model="latest">
              {{$L('latest result of each job on each node')}}
            </el-checkbox>
          </el-col>
          <el-col :sm="12" :xs="24">
            <el-checkbox v-model="failedOnly">
              {{$L('failure only')}}
            </el-checkbox>
          </el-col>
        </el-row>
        <div class="log-filter__actions">
          <el-button type="primary" @click="submit">{{$L('submit query')}}</el-button>
        </div>
      </el-form>
    </el-card>

    <el-table :data="list" border stripe v-loading="loading" :empty-text="$L('no data')">
      <el-table-column :label="$L('job name')" min-width="200">
        <template slot-scope="{row}">
          <router-link :to="'/job/edit/'+row.jobGroup+'/'+row.jobId">{{row.name}}</router-link>
        </template>
      </el-table-column>
      <el-table-column :label="$L('executing node')" min-width="200">
        <template slot-scope="{row}">
          <span :title="row.node">{{$store.getters.hostshows(row.node)}}</span>
        </template>
      </el-table-column>
      <el-table-column prop="user" :label="$L('executing user')" min-width="160"></el-table-column>
      <el-table-column :label="$L('executing time')" min-width="240">
        <template slot-scope="{row}">
          <span :class="{'log-warning': durationAttention(row.beginTime, row.endTime)}">
            {{formatTime(row)}}
          </span>
        </template>
      </el-table-column>
      <el-table-column :label="$L('executing result')" min-width="220">
        <template slot-scope="{row}">
          <router-link :to="'/log/'+row.id">
            {{$L(row.success ? 'successed' : 'failed')}}
          </router-link>
          |
          <el-tooltip :content="$L('click to select a node and re-execute job')">
            <el-button type="text" icon="el-icon-refresh" @click.prevent="showExecuteJobModal(row.name, row.jobGroup, row.jobId)"></el-button>
          </el-tooltip>
        </template>
      </el-table-column>
    </el-table>
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

      // init default date range to today to avoid heavy queries
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
