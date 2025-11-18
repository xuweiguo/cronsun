<style scoped>
.job-rule {
  margin-bottom: 16px;
}

.job-rule .layui-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.job-rule .layui-form-item {
  margin-bottom: 16px;
}

.rule-help {
  font-size: 12px;
  color: #999;
  margin-top: 4px;
}
</style>

<template>
<div class="layui-card job-rule">
  <div class="layui-card-header">
    <span>{{$L('timer')}} - {{index}}</span>
    <button class="layui-btn layui-btn-danger layui-btn-xs" @click.prevent="remove">{{$L('delete')}}</button>
  </div>
  <div class="layui-card-body">
    <div class="layui-form-item">
      <label class="layui-form-label">{{$L('cron expression')}}</label>
      <div class="layui-input-block">
        <input type="text" class="layui-input" :value="rule.timer" @input="change('timer', $event.target.value)" :placeholder="$L('0 * * * * *, rules see the 「?」on the right')"/>
        <div class="rule-help">{{$L('<sec> <min> <hr> <day> <month> <week>, rules is same with Cron')}}</div>
      </div>
    </div>
    <div class="layui-form-item">
      <label class="layui-form-label">{{$L('node group')}}</label>
      <div class="layui-input-block">
        <Dropdown :title="$L('node group')" :items="nodeGroups" :selected="rule.gids" multiple="true" @change="changeNodeGroups" />
      </div>
    </div>
    <div class="layui-form-item">
      <label class="layui-form-label">{{$L('and please running on those nodes')}}</label>
      <div class="layui-input-block">
        <Dropdown :title="$L('select nodes')" :items="activityNodes" :selected="rule.nids" multiple="true" @change="changeIncludeNodes" />
      </div>
    </div>
    <div class="layui-form-item">
      <label class="layui-form-label">{{$L('do not running on those nodes')}}</label>
      <div class="layui-input-block">
        <Dropdown :title="$L('select nodes')" :items="activityNodes" :selected="rule.exclude_nids" multiple="true" @change="changeExcludeNodes" />
      </div>
    </div>
  </div>
</div>
</template>

<script>
import Dropdown from './basic/Dropdown.vue';

export default {
  name: 'job-edit-rule',
  props: ['index', 'rule'],
  data: function(){
    return {
      nodeGroups: [],
      activityNodes: []
    }
  },

  mounted: function(){
    var vm = this;
    this.activityNodes = this.$store.getters.dropdownNodes;

    this.$rest.GET('node/groups').onsucceed(200, (resp)=>{
      var groups = [];
      for (var i in resp) {
        groups.push({value: resp[i].id, name: resp[i].name});
      }
      vm.nodeGroups = groups;
    }).do();
  },

  methods: {
    remove: function(){
      this.$emit('remove', this.index);
    },
    change: function(key, val){
      this.$emit('change', this.index, key, val);
    },
    changeNodeGroups: function(val){
      var gids = Array.isArray(val) ? val : (val && val.trim().length === 0 ? [] : (val || '').split(','));
      this.change('gids', gids);
    },
    changeIncludeNodes: function(val){
      var nids = Array.isArray(val) ? val : (val && val.trim().length === 0 ? [] : (val || '').split(','));
      this.change('nids', nids);
    },
    changeExcludeNodes: function(val){
      var exclude_nids = Array.isArray(val) ? val : (val && val.trim().length === 0 ? [] : (val || '').split(','));
      this.change('exclude_nids', exclude_nids);
    }
  },

  components: {
    Dropdown
  }
}
</script>
