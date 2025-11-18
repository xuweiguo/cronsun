<template>
  <div class="cs-pager" v-if="total > 0">
    <div class="layui-btn-group" v-if="total > 1">
      <button class="layui-btn layui-btn-sm" :disabled="_current <= 1" @click="go(_current - 1)">&laquo;</button>
      <button
        v-for="page in pageNumbers"
        :key="page"
        class="layui-btn layui-btn-sm"
        :class="{'layui-btn-normal': page === _current}"
        @click="go(page)">
        {{page}}
      </button>
      <button class="layui-btn layui-btn-sm" :disabled="_current >= total" @click="go(_current + 1)">&raquo;</button>
    </div>
    <div class="cs-pager__summary">{{_current}} / {{total}}</div>
  </div>
</template>

<script>
export default {
  name: 'pager',
  props: {
    total: {type: Number, default: 1, required: true},
    maxBtn: {type: Number, default: 7},
    pageVar: {type: String, default: 'page'}
  },
  data: function(){
    return {
      _pagevar: 'page',
      _current: 1
    }
  },

  created: function(){
    this._pagevar = this.pageVar || 'page';
    this._current = parseInt(this.$route.query[this._pagevar]) || 1;
  },

  methods: {
    pageURL: function(page){
      return this.url + this._pagevar + '=' + page;
    },

    go: function(page){
      if (page < 1 || page > this.total || page === this._current) return;
      this.$router.push(this.pageURL(page));
    }
  },

  watch: {
    '$route': function(){
      this._current = parseInt(this.$route.query[this._pagevar]) || 1;
    }
  },

  computed: {
    url: function(){
      var query = [];
      for (var  k in this.$route.query) {
        if (this._pagevar === k) continue;
        query.push(k+'='+this.$route.query[k]);
      }

      var queryStr = query.join('&');
      return queryStr.length ? this.$route.path+'?'+queryStr + '&' : this.$route.path+'?';
    },

    pageNumbers() {
      if (this.total <= this.maxBtn) {
        return Array.from({length: this.total}, (_, i) => i + 1);
      }
      const half = Math.floor(this.maxBtn / 2);
      let start = this._current - half;
      let end = this._current + half;
      if (start < 1) {
        start = 1;
        end = this.maxBtn;
      }
      if (end > this.total) {
        end = this.total;
        start = this.total - this.maxBtn + 1;
      }
      const pages = [];
      for (let i = start; i <= end; i += 1) {
        pages.push(i);
      }
      return pages;
    }
  }
}
</script>

<style scoped>
.cs-pager {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin: 16px 0;
  flex-wrap: wrap;
  gap: 12px;
}

.cs-pager__summary {
  color: #606266;
  font-size: 14px;
  min-width: 70px;
  text-align: right;
}
</style>
