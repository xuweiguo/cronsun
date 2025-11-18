<style scoped>
.execute-modal {
  width: 420px;
  left: 50% !important;
  margin-left: -210px;
}

.execute-modal .layui-layer-content {
  padding: 16px;
}

.modal-buttons {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 12px 16px 16px;
}
</style>

<template>
  <div v-if="visible">
    <div class="layui-layer-shade" @click="hide"></div>
    <div class="layui-layer layui-layer-page execute-modal">
      <div class="layui-layer-title">{{$L('executing job: {job}', jobName)}}</div>
      <div class="layui-layer-content">
        <Dropdown
          :title="$L('node')"
          :items="nodes"
          v-on:change="changeNode"
          style="width:100%"/>
      </div>
      <div class="modal-buttons">
        <button class="layui-btn layui-btn-primary" type="button" @click="hide">{{$L('cancel')}}</button>
        <button class="layui-btn" type="button" :disabled="loading" @click="submit">{{$L('execute now')}}</button>
      </div>
    </div>
  </div>
</template>

<script>
import Dropdown from './basic/Dropdown.vue';

export default {
  name: 'execute-job',
  data(){
    return {
      jobGroup: '',
      jobId: '',
      jobName: '',
      nodes: [],
      selectedNode: '',
      loading: false,
      visible: false
    }
  },

  methods: {
    show(jobName, jobGroup, jobId){
      this.jobName = jobName;
      this.jobGroup = jobGroup;
      this.jobId = jobId;
      this.visible = true;
      this.fetchJobNodes();
    },

    hide(){
      this.visible = false;
      this.loading = false;
    },

    fetchJobNodes(){
      var vm = this;
      this.loading = true;
      this.$rest.GET('job/'+this.jobGroup+'-'+this.jobId+'/nodes').
        onsucceed(200, (resp)=>{
          var nodes = [{value: 'all nodes', name: vm.$L('all nodes')}];
          for (var i in resp) {
            nodes.push({value: resp[i], name: vm.$store.getters.hostshows(resp[i])})
          }
          vm.nodes = nodes;
        }).
        onfailed((msg)=>{
          vm.$bus.$emit('error', msg);
          vm.hide();
        }).
        onend(()=>{vm.loading = false}).
        do();
    },

    submit(){
      if (this.loading) return;
      var vm = this;
      this.loading = true;
      var node = this.selectedNode === 'all nodes' ? '' : this.selectedNode;
      this.$rest.PUT('/job/'+this.jobGroup+'-'+this.jobId+'/execute?node='+node).
        onsucceed(204, ()=>{
          vm.$bus.$emit('success', '执行命令已发送，注意查看任务日志');
          vm.hide();
        }).
        onfailed((msg)=>{vm.$bus.$emit('error', msg)}).
        onend(()=>{vm.loading = false}).
        do();
      return false;
    },

    changeNode(val){
      this.selectedNode = val;
    }
  },

  components: {
    Dropdown
  }
}
</script>
