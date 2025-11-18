<style scoped>
.job-edit-form {
  background: #fff;
  padding: 24px;
  border-radius: 4px;
  position: relative;
}

.job-loading {
  position: absolute;
  inset: 0;
  background: rgba(255,255,255,0.8);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  z-index: 2;
}

.job-edit-header {
  display: flex;
  align-items: center;
  gap: 16px;
  flex-wrap: wrap;
  margin-bottom: 16px;
}

.job-edit-header h3 {
  margin: 0;
}

.job-toggle-group {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.job-type-options,
.job-form-actions {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.job-form-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 16px;
}

.job-form-item {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.job-alert {
  padding: 12px;
  border: 1px dashed #f8d7da;
  background: #fff5f5;
  color: #d93025;
  margin-bottom: 16px;
}
</style>

<template>
  <div class="job-edit-form">
    <div class="job-loading" v-if="loading">{{$L('loading')}}...</div>
    <form class="layui-form" @submit.prevent="submit">
      <div class="job-edit-header">
        <h3>{{$L(action == 'CREATE' ? 'create job' : 'update job')}} <em v-if="job.id">ID# {{job.id}}</em></h3>
        <div class="job-toggle-group">
          <label><input type="checkbox" v-model="job.pause"> {{$L('pause')}}</label>
          <label v-if="$appConfig.alarm"><input type="checkbox" v-model="job.fail_notify"> {{job.fail_notify ? $L('warning on') : $L('warning off')}}</label>
        </div>
      </div>

      <div class="job-form-item">
        <label>{{$L('job type')}}</label>
        <div class="job-type-options">
          <label><input type="radio" :value="0" v-model.number="job.kind"> {{$L('common job')}}</label>
          <label><input type="radio" :value="1" v-model.number="job.kind"> {{$L('single node single process')}}</label>
          <label><input type="radio" :value="2" v-model.number="job.kind"> {{$L('group level common')}}</label>
        </div>
      </div>

      <div class="job-form-grid">
        <div class="job-form-item">
          <label>{{$L('job name')}}</label>
          <input class="layui-input" :value="job.name" @input="updateValue" :placeholder="$L('job name')">
        </div>
        <div class="job-form-item">
          <label>{{$L('job group')}}</label>
          <Dropdown :title="$L('select a group')" :allowAdditions="true" :items="groups" :selected="job.group" @change="changeGroup"/>
        </div>
      </div>

      <div class="job-form-grid">
        <div class="job-form-item" style="grid-column: span 2;">
          <label>{{$L('script path')}} {{allowSuffixsTip}}</label>
          <input class="layui-input" v-model="job.cmd" :placeholder="$L('script path')">
        </div>
        <div class="job-form-item">
          <label>{{$L($appConfig.security.open ? 'user(required)' : 'user(optional)')}} </label>
          <Dropdown
            v-if="$appConfig.security.open"
            :title="$L('the user which to execute the command')"
            :items="$appConfig.security.users"
            :selected="job.user"
            @change="changeUser"/>
          <input
            v-else
            class="layui-input"
            v-model="job.user"
            :placeholder="$L('the user which to execute the command')">
        </div>
      </div>

      <div class="job-form-item" v-if="$appConfig.alarm && job.fail_notify">
        <label>{{$L('warning receiver')}}</label>
        <Dropdown :title="$L('e-mail address')" :items="alarmReceivers" :selected="job.to" @change="changeAlarmReceiver" :multiple="true" :allowAdditions="true"/>
      </div>

      <div class="job-form-grid">
        <div class="job-form-item">
          <label>{{$L('timeout(in seconds, 0 for no limits)')}}</label>
          <input class="layui-input" type="number" min="0" v-model.number="job.timeout">
        </div>
        <div class="job-form-item" v-show="job.kind === 0">
          <label>{{$L('parallel number in one node(0 for no limits)')}}</label>
          <input class="layui-input" type="number" min="0" v-model.number="job.parallels">
        </div>
      </div>

      <div class="job-form-grid">
        <div class="job-form-item">
          <label>{{$L('retries(number of retries when failed, 0 means no retry)')}}</label>
          <input class="layui-input" type="number" min="0" v-model.number="job.retry">
        </div>
        <div class="job-form-item">
          <label>{{$L('retry interval(in seconds)')}}</label>
          <input class="layui-input" type="number" min="0" v-model.number="job.interval">
        </div>
      </div>

      <div class="job-form-item" v-if="$appConfig.log_expiration_days>0">
        <label>{{$L('log expiration(log expired after N days, 0 will use default setting: {n} days)', $appConfig.log_expiration_days)}} </label>
        <input class="layui-input" type="number" min="0" v-model.number="job.log_expiration">
      </div>

      <div class="job-alert" v-if="!job.rules || job.rules.length == 0">
        {{$L('the job dose not have a timer currently, please click the button below to add a timer')}}
      </div>

      <JobEditRule v-for="(rule, index) in job.rules" :key="rule.id" :rule="rule" :index="index" @remove="removeRule" @change="changeRule"/>

      <div class="job-form-actions">
        <button class="layui-btn layui-btn-primary" type="button" @click="addNewTimer">{{$L('add timer')}}</button>
        <button class="layui-btn" type="submit">{{$L('save job')}}</button>
      </div>
    </form>
  </div>
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
