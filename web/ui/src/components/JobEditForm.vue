<style scoped>
.job-edit-header {
  display: flex;
  align-items: center;
  gap: 16px;
  flex-wrap: wrap;
}

.job-edit-actions {
  display: flex;
  gap: 12px;
}

.job-form-actions {
  display: flex;
  gap: 12px;
}
</style>

<template>
  <el-form class="job-edit-form" label-position="top" v-loading="loading">
    <div class="job-edit-header">
      <h3>{{$L(action == 'CREATE' ? 'create job' : 'update job')}} <em v-if="job.id">ID# {{job.id}}</em></h3>
      <el-switch
        v-model="job.pause"
        :active-value="true"
        :inactive-value="false"
        :active-text="$L('pause')"
        :inactive-text="$L('open')"
        active-color="#F56C6C"
        inactive-color="#67C23A">
      </el-switch>
      <el-switch
        v-if="$appConfig.alarm"
        v-model="job.fail_notify"
        :active-value="true"
        :inactive-value="false"
        :active-text="$L('warning on')"
        :inactive-text="$L('warning off')">
      </el-switch>
    </div>

    <el-form-item :label="$L('job type')">
      <el-radio-group v-model="job.kind">
        <el-radio-button :label="0">{{$L('common job')}}</el-radio-button>
        <el-radio-button :label="1">{{$L('single node single process')}}</el-radio-button>
        <el-radio-button :label="2">
          <el-tooltip :content="$L('group level common help')">
            <span>{{$L('group level common')}}</span>
          </el-tooltip>
        </el-radio-button>
      </el-radio-group>
    </el-form-item>

    <el-row :gutter="16">
      <el-col :sm="12" :xs="24">
        <el-form-item :label="$L('job name')">
          <el-input :value="job.name" @input="updateValue" :placeholder="$L('job name')"/>
        </el-form-item>
      </el-col>
      <el-col :sm="12" :xs="24">
        <el-form-item :label="$L('job group')">
          <Dropdown :title="$L('select a group')" :allowAdditions="true" :items="groups" :selected="job.group" @change="changeGroup"/>
        </el-form-item>
      </el-col>
    </el-row>

    <el-row :gutter="16">
      <el-col :sm="18" :xs="24">
        <el-form-item :label="$L('script path') + ' ' + allowSuffixsTip">
          <el-input v-model="job.cmd" :placeholder="$L('script path')"/>
        </el-form-item>
      </el-col>
      <el-col :sm="6" :xs="24">
        <el-form-item :label="$L($appConfig.security.open ? 'user(required)' : 'user(optional)')">
          <Dropdown
            v-if="$appConfig.security.open"
            :title="$L('the user which to execute the command')"
            :items="$appConfig.security.users"
            :selected="job.user"
            @change="changeUser"/>
          <el-input
            v-else
            v-model="job.user"
            :placeholder="$L('the user which to execute the command')"/>
        </el-form-item>
      </el-col>
    </el-row>

    <el-form-item v-if="$appConfig.alarm && job.fail_notify" :label="$L('warning receiver')">
      <Dropdown :title="$L('e-mail address')" :items="alarmReceivers" :selected="job.to" @change="changeAlarmReceiver" :multiple="true" :allowAdditions="true"/>
    </el-form-item>

    <el-row :gutter="16">
      <el-col :sm="12" :xs="24">
        <el-form-item :label="$L('timeout(in seconds, 0 for no limits)')">
          <el-input-number v-model.number="job.timeout" :min="0" :controls="false" style="width:100%"/>
        </el-form-item>
      </el-col>
      <el-col :sm="12" :xs="24" v-show="job.kind === 0">
        <el-form-item :label="$L('parallel number in one node(0 for no limits)')">
          <el-input-number v-model.number="job.parallels" :min="0" :controls="false" style="width:100%"/>
        </el-form-item>
      </el-col>
    </el-row>

    <el-row :gutter="16">
      <el-col :sm="12" :xs="24">
        <el-form-item :label="$L('retries(number of retries when failed, 0 means no retry)')">
          <el-input-number v-model.number="job.retry" :min="0" :controls="false" style="width:100%"/>
        </el-form-item>
      </el-col>
      <el-col :sm="12" :xs="24">
        <el-form-item :label="$L('retry interval(in seconds)')">
          <el-input-number v-model.number="job.interval" :min="0" :controls="false" style="width:100%"/>
        </el-form-item>
      </el-col>
    </el-row>

    <el-form-item v-if="$appConfig.log_expiration_days>0" :label="$L('log expiration(log expired after N days, 0 will use default setting: {n} days)', $appConfig.log_expiration_days)">
      <el-input-number v-model.number="job.log_expiration" :min="0" :controls="false" style="width:100%"/>
    </el-form-item>

    <div class="field">
      <el-alert v-if="!job.rules || job.rules.length == 0" type="warning" :title="$L('the job dose not have a timer currently, please click the button below to add a timer')" show-icon></el-alert>
    </div>

    <JobEditRule v-for="(rule, index) in job.rules" :key="rule.id" :rule="rule" :index="index" @remove="removeRule" @change="changeRule"/>

    <div class="job-form-actions">
      <el-button type="info" icon="el-icon-time" @click="addNewTimer">{{$L('add timer')}}</el-button>
      <el-button type="primary" icon="el-icon-upload" @click="submit">{{$L('save job')}}</el-button>
    </div>
  </el-form>
</template>

<script>
import JobEditRule from './JobEditRule.vue';
import Dropdown from './basic/Dropdown.vue';
import {split} from '../libraries/functions';

export default {
  name: 'job-edit-form',
  data: function(){
      return {
        action: 'CREATE',
        groups: [],
        alarmReceivers: [],
        loading: false,
        allowSuffixsTip: '',
        job: {
          id: '',
          name: '',
          group: '',
          oldGroup: '',
          cmd: '',
          user: '',
          kind: 0,
          pause: false,
          parallels: 0,
          timeout: 0,
          interval: 0,
          retry: 0,
          rules: [],
          fail_notify: false,
          log_expiration: 0,
          to: []
        }
      }
  },

  methods: {
    updateValue: function(v){
      var tv = v.replace(/[\*\/]/g, '');
      this.job.name = tv;
    },

    addNewTimer: function(){
      if (!this.job.rules) this.job.rules = [];
      this.job.rules.push({id: this.newRandomRuleId()});
    },

    changeGroup: function(val){
      this.job.group = val;
    },

    changeUser: function(val){
      this.job.user = val;
    },

    changeAlarmReceiver: function(val){
      this.job.to = Array.isArray(val) ? val : split(val, ',');
    },

    removeRule: function(index){
      this.job.rules.splice(index, 1);
    },

    changeRule: function(index, key, val){
      this.job.rules[index][key] = val;
    },

    submit: function(){
      var exceptCode = this.action == 'CREATE' ? 201 : 200;
      this.loading = true;
      this.$rest.PUT('job', this.job)
        .onsucceed(exceptCode, ()=>{this.$router.push({path:'/job',query:this.$route.query})})
        .onfailed((resp)=>{this.$bus.$emit('error', resp)})
        .onend(()=>{this.loading=false})
        .do();
    },

    newRandomRuleId: function(){
      return 'NEW'+Math.random().toString();
    }
  },

  mounted: function(){
    var secCnf = this.$appConfig.security;
    if (secCnf.open && secCnf.ext && secCnf.ext.length > 0) {
      this.allowSuffixsTip = this.$L('(only [{.suffixs}] files can be allowed)', secCnf.ext.join(' '));
    }

    if (this.$route.path.indexOf('/job/create') === 0) {
      this.action = 'CREATE';
    } else {
      this.action = 'UPDATE';
      this.loading = true;
      this.$rest.GET('job/'+this.$route.params.group+'-'+this.$route.params.id).
        onsucceed(200, (resp)=>{
          this.job = resp;
          this.alarmReceivers = resp.to;
          this.job.oldGroup = resp.group;
          if (this.job.rules) {
            for (var i in this.job.rules) {
              if (this.job.rules[i].id.length == 0) {
                this.job.rules[i].id = this.newRandomRuleId();
              }
            }
          }
        }).
        onfailed((msg)=> this.$bus.$emit('error', msg)).
        onend(()=>{this.loading = false}).
        do();
    }

    this.$rest.GET('job/groups').onsucceed(200, (resp)=>{
      !resp.includes('default') && resp.unshift('default');
      this.groups = resp;
    }).do();
  },

  components: {
    JobEditRule,
    Dropdown
  }
}
</script>
