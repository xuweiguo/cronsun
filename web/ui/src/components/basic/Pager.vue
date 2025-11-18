<template>
  <div class="cs-pager">
    <el-pagination
      v-if="total > 1"
      background
      layout="prev, pager, next, jumper"
      :pager-count="maxBtn"
      :current-page="_current"
      :page-count="total"
      @current-change="go"/>
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
      if (page < 1 || page > this.total) return;
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
}

.cs-pager__summary {
  color: #606266;
  font-size: 14px;
  min-width: 70px;
  text-align: right;
}
</style>
