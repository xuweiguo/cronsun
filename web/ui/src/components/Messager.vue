<style scoped>
.message-container {
  position: fixed;
  top: 80px;
  right: 24px;
  width: 360px;
  z-index: 3000;
}

.message-item {
  margin-bottom: 12px;
}
</style>

<template>
  <div class="message-container">
    <transition-group name="el-fade-in-linear" tag="div">
      <el-alert
        v-for="item in queue"
        :key="item.id"
        :title="item.content"
        :type="item.type"
        show-icon
        closable
        class="message-item"
        @close="closeMessage(item.id)">
      </el-alert>
    </transition-group>
  </div>
</template>

<script>
export default {
  name: 'message',
  data(){
    return {
      queue: []
    }
  },

  methods: {
    showMessage(type, content){
      var id = Math.random().toString();
      this.queue.push({
        id: id,
        content: content,
        type: type
      });

      var vm = this;
      setTimeout(()=>{
        vm.closeMessage(id);
      }, 5000);
    },

    closeMessage(id){
      var index = this.queue.findIndex((m) => m.id === id);
      if (index !== -1) {
        this.queue.splice(index, 1);
      }
    }
  },

  mounted(){
    var vm = this;
    this.$bus.$on('error', (content)=>{
      vm.showMessage('error', content);
    });
    this.$bus.$on('success', (content)=>{
      vm.showMessage('success', content);
    });
    this.$bus.$on('warning', (content)=>{
      vm.showMessage('warning', content);
    });
  }
}
</script>