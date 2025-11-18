<style scoped>
.app-layout {
  min-height: 100vh;
  background: #f5f7fa;
}

.app-header {
  background: #1f2d3d;
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
}

.app-header__left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.app-logo {
  font-size: 20px;
  font-weight: 600;
  letter-spacing: 0.1em;
}

.app-menu {
  border-bottom: none;
}

.app-header__actions {
  display: flex;
  align-items: center;
  gap: 12px;
  color: #fff;
}

.header-link {
  color: #fff;
  display: flex;
  align-items: center;
  gap: 4px;
}

.header-link:hover {
  color: #ffd04b;
}

.app-main {
  padding: 24px;
}

.setting-section {
  margin-bottom: 24px;
}

.setting-section__title {
  font-size: 14px;
  color: #909399;
  margin-bottom: 8px;
}

.drawer-version {
  font-weight: 600;
  margin-bottom: 16px;
}

.drawer-group {
  margin-bottom: 16px;
}

.drawer-group:last-child {
  margin-bottom: 0;
}

.drawer-menu {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.drawer-menu .el-radio-button {
  margin-bottom: 8px;
}

.app-header__actions .el-button.is-circle {
  color: #fff;
}
</style>

<template>
  <div id="app">
    <el-container class="app-layout">
      <el-header class="app-header">
        <div class="app-header__left">
          <span class="app-logo">CRONSUN</span>
          <el-menu
            class="app-menu"
            mode="horizontal"
            background-color="#1f2d3d"
            text-color="#bfcbd9"
            active-text-color="#ffd04b"
            :default-active="activeMenu"
            @select="navigate">
            <el-menu-item v-if="shouldOpen" index="/">
              <i class="el-icon-menu"></i>
              {{$L('dashboard')}}
            </el-menu-item>
            <el-menu-item v-if="shouldOpen" index="/log">
              <i class="el-icon-document"></i>
              {{$L('log')}}
            </el-menu-item>
            <el-menu-item v-if="shouldOpen" index="/job">
              <i class="el-icon-date"></i>
              {{$L('job')}}
            </el-menu-item>
            <el-menu-item v-if="shouldOpen" index="/node">
              <i class="el-icon-s-platform"></i>
              {{$L('node')}}
            </el-menu-item>
            <el-menu-item
              v-if="$store.getters.enabledAuth && $store.getters.role === 1"
              index="/admin/account/list">
              <i class="el-icon-user-solid"></i>
              {{$L('account')}}
            </el-menu-item>
            <el-menu-item index="docs" @click="openDocs">
              <i class="el-icon-link"></i>
              Docs
            </el-menu-item>
          </el-menu>
        </div>
        <div class="app-header__actions">
          <router-link
            v-if="$store.getters.email"
            to="/user/setpwd"
            class="header-link">
            <i class="el-icon-user"></i>
            {{$store.getters.email}}
          </router-link>
          <el-button
            v-if="$store.getters.email"
            type="text"
            circle
            icon="el-icon-switch-button"
            @click="logout"></el-button>
          <el-button
            type="text"
            circle
            icon="el-icon-setting"
            @click="toggleSetting"></el-button>
        </div>
      </el-header>
      <el-main class="app-main">
        <router-view></router-view>
      </el-main>
    </el-container>

    <el-drawer
      title="{{$L('settings')}}"
      :visible.sync="settingOpen"
      direction="rtl"
      size="320px">
      <div class="drawer-version">{{$store.getters.version}}</div>
      <div class="setting-section">
        <div class="setting-section__title">Language</div>
        <el-radio-group v-model="locale" size="small" @change="selectLanguage">
          <el-radio-button
            v-for="lang in $Lang.supported"
            :key="lang.code"
            :label="lang.code">
            {{lang.name}}
          </el-radio-button>
        </el-radio-group>
      </div>
      <div class="setting-section">
        <div class="setting-section__title">{{$L('node show as')}}</div>
        <el-radio-group v-model="hostshow" size="small" @change="selectHostshows">
          <el-radio-button
            v-for="h in hostshowsList"
            :key="h.value"
            :label="h.value">
            {{h.name}}
          </el-radio-button>
        </el-radio-group>
      </div>
    </el-drawer>

    <Messager/>
  </div>
</template>

<script>
import store from './vuex/store';
import Vue from 'vue';
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
      if (index === 'docs') {
        this.openDocs();
        return;
      }
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
