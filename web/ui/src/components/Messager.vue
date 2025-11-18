<template>
  <div></div>
</template>

<script>
export default {
  name: 'message',
  data(){
    return {
      layer: null
    }
  },

  methods: {
    showMessage(type, content){
      const iconMap = { success: 1, error: 2, warning: 0 };
      if (this.layer && typeof this.layer.msg === 'function') {
        this.layer.msg(content, {icon: iconMap[type] || 0, time: 3000});
        return;
      }
      window.alert(content);
    }
  },

  mounted(){
    const layui = this.$layui || window.layui || {};
    this.layer = layui.layer;
    this.$bus.$on('error', (content)=>{
      this.showMessage('error', content);
    });
    this.$bus.$on('success', (content)=>{
      this.showMessage('success', content);
    });
    this.$bus.$on('warning', (content)=>{
      this.showMessage('warning', content);
    });
  }
}
</script>
