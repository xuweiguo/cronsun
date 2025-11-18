<style scoped>
.app-layout {
  min-height: 100vh;
  background: #f2f2f2;
}

.app-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-right: 16px;
}

.app-header__left {
  display: flex;
  align-items: center;
  gap: 24px;
}

.app-logo {
  font-size: 20px;
  font-weight: 700;
  letter-spacing: 0.08em;
  color: #fff;
}

.app-menu {
  background: transparent;
}

.app-menu .layui-nav-item {
  line-height: 60px;
}

.app-header__actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.header-link {
  color: #fff;
}

.header-link:hover {
  text-decoration: underline;
}

.app-main {
  padding: 24px;
  background: #f2f2f2;
}

.setting-drawer {
  position: fixed;
  top: 0;
  right: -360px;
  width: 320px;
  height: 100vh;
  background: #fff;
  box-shadow: -2px 0 8px rgba(0, 0, 0, 0.15);
  transition: right 0.2s ease-in-out;
  z-index: 1000;
  padding: 24px;
}

.setting-drawer.is-open {
  right: 0;
}

.setting-mask {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.35);
  z-index: 999;
}

.setting-section {
  margin-bottom: 24px;
}

.setting-section__title {
  font-size: 13px;
  color: #666;
  margin-bottom: 8px;
  font-weight: 600;
}

.setting-radio-group {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.setting-radio-group label {
  cursor: pointer;
  padding: 4px 12px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
}

.setting-radio-group input {
  margin-right: 4px;
}
</style>

<template>
  <div id="app">
    <div class="layui-layout layui-layout-admin app-layout">
      <div class="layui-header app-header">
        <div class="app-header__left">
          <span class="app-logo">CRONSUN</span>
          <ul class="layui-nav layui-bg-black app-menu">
            <li
              v-if="shouldOpen"
              class="layui-nav-item"
              :class="{'layui-this': activeMenu === '/'}"
              @click="navigate('/')">
              <a>{{$L('dashboard')}}</a>
            </li>
            <li
              v-if="shouldOpen"
              class="layui-nav-item"
              :class="{'layui-this': activeMenu === '/log'}"
              @click="navigate('/log')">
              <a>{{$L('log')}}</a>
            </li>
            <li
              v-if="shouldOpen"
              class="layui-nav-item"
              :class="{'layui-this': activeMenu === '/job'}"
              @click="navigate('/job')">
              <a>{{$L('job')}}</a>
            </li>
            <li
              v-if="shouldOpen"
              class="layui-nav-item"
              :class="{'layui-this': activeMenu === '/node'}"
              @click="navigate('/node')">
              <a>{{$L('node')}}</a>
            </li>
            <li
              v-if="$store.getters.enabledAuth && $store.getters.role === 1"
              class="layui-nav-item"
              :class="{'layui-this': activeMenu === '/admin/account/list'}"
              @click="navigate('/admin/account/list')">
              <a>{{$L('account')}}</a>
            </li>
            <li class="layui-nav-item" @click="openDocs">
              <a>Docs</a>
            </li>
          </ul>
        </div>
        <div class="app-header__actions">
          <router-link
            v-if="$store.getters.email"
            to="/user/setpwd"
            class="header-link">
            {{$store.getters.email}}
          </router-link>
          <button
            v-if="$store.getters.email"
            class="layui-btn layui-btn-sm layui-btn-primary"
            type="button"
            @click="logout">
            {{$L('logout')}}
          </button>
          <button
            class="layui-btn layui-btn-sm"
            type="button"
            @click="toggleSetting">
            {{$L('settings')}}
          </button>
        </div>
      </div>
      <div class="layui-body app-main">
        <router-view></router-view>
      </div>
    </div>

    <div class="setting-mask" v-if="settingOpen" @click="toggleSetting"></div>
    <div class="setting-drawer" :class="{'is-open': settingOpen}">
      <h3>{{$L('settings')}}</h3>
      <div class="setting-section">
        <div class="setting-section__title">Version</div>
        <div>{{$store.getters.version}}</div>
      </div>
      <div class="setting-section">
        <div class="setting-section__title">Language</div>
        <div class="setting-radio-group">
          <label v-for="lang in $Lang.supported" :key="lang.code">
            <input type="radio" name="lang" :value="lang.code" v-model="locale" @change="selectLanguage(lang.code)">
            {{lang.name}}
          </label>
        </div>
      </div>
      <div class="setting-section">
        <div class="setting-section__title">{{$L('node show as')}}</div>
        <div class="setting-radio-group">
          <label v-for="h in hostshowsList" :key="h.value">
            <input type="radio" name="hostshows" :value="h.value" v-model="hostshow" @change="selectHostshows(h.value)">
            {{h.name}}
          </label>
        </div>
      </div>
    </div>

    <Messager/>
  </div>
</template>

<script>
import store from './vuex/store';
import Cookies from 'js-cookie';
import Messager from './components/Messager.vue';

export default {
  name: 'app',
  store,

  data: function() {
    return {
      hostshowsList: [{name: this.$L('hostname'), value: 'hostname'}, {name: 'IP', value: 'ip'}],
      locale: '',
      settingOpen: false
    }
  },

  created: function() {
    this.locale = Cookies.get('locale') || this.$Lang.getLocale();
    this.$store.commit('setShowWithHostname', Cookies.get('hostshows') === 'hostname');
  },

  computed: {
    shouldOpen() {
      return !this.$store.getters.enabledAuth || (this.$store.getters.enabledAuth && this.$store.getters.email)
    },

    hostshow: {
      get() {
        return this.$store.getters.showWithHostname ? 'hostname' : 'ip';
      },
      set(value) {
        const showHost = value === 'hostname';
        this.$store.commit('setShowWithHostname', showHost);
      }
    },

    activeMenu() {
      const path = this.$route.path;
      if (path.indexOf('/log') === 0) return '/log';
      if (path.indexOf('/job') === 0) return '/job';
      if (path.indexOf('/node') === 0) return '/node';
      if (path.indexOf('/admin/account') === 0) return '/admin/account/list';
      if (path.indexOf('/user') === 0) return '/user/setpwd';
      return '/';
    }
  },

  methods: {
    logout() {
      var vm = this;
      this.$rest.DELETE('session').
        onsucceed(200, ()=>{
          vm.$store.commit('setEmail', '');
          vm.$store.commit('setRole', 0);
          vm.$router.push('/login');
        }).
        do();
    },

    toggleSetting() {
      this.settingOpen = !this.settingOpen;
    },

    openDocs() {
      window.open('https://github.com/shunfei/cronsun/wiki', '_blank');
    },

    navigate(index) {
      if (index && this.$route.path !== index) {
        this.$router.push(index);
      }
    },

    selectLanguage(code) {
      if (this.locale !== code) {
        Cookies.set('locale', code);
        window.location.reload();
        return;
      }
      this.settingOpen = false;
    },

    selectHostshows(v) {
      if (this.hostshow !== v) {
        Cookies.set('hostshows', v);
        window.location.reload();
        return;
      }
      this.settingOpen = false;
    }
  },

  components: {
    Messager
  },
}
</script>
