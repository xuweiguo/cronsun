<template>
  <el-dialog
    :title="$L('executing job: {job}', jobName)"
    :visible.sync="visible"
    width="480px"
    :close-on-click-modal="false">
    <Dropdown
      :title="$L('node')"
      :items="nodes"
      v-on:change="changeNode"
      style="width:100%"/>
    <span slot="footer" class="dialog-footer">
      <el-button @click="hide">{{$L('cancel')}}</el-button>
      <el-button type="primary" :loading="loading" @click="submit">{{$L('execute now')}}</el-button>
    </span>
  </el-dialog>
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
