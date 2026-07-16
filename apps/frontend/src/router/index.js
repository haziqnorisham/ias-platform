import { createRouter, createWebHistory } from 'vue-router';
import { useAuth } from '../composables/useAuth';

// route components
import Home from '../views/Home.vue';
import About from '../views/About.vue';
import Devices from '../views/Devices.vue';
import Diagnostics from '../views/Diagnostics.vue';
import Ai from '../views/Ai.vue';
import Dashboards from '../views/Dashboards.vue';
import DashboardEdit from '../views/DashboardEdit.vue';
import DeviceProfiles from '../views/DeviceProfiles.vue';
import DataBrowser from '../views/DataBrowser.vue';
import IngestLogs from '../views/IngestLogs.vue';
import Settings from '../views/Settings.vue';
import ExtensionsHub from '../views/ExtensionsHub.vue';
import ExtensionDetail from '../views/ExtensionDetail.vue';
import DashboardView from '../views/DashboardView.vue';
import Login from '../views/Login.vue';

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home
  },
  {
    path: '/about',
    name: 'About',
    component: About
  },
  {
    path: '/devices',
    name: 'Devices',
    component: Devices
  },
  {
    path: '/diagnostics',
    name: 'Diagnostics',
    component: Diagnostics
  },
  {
    path: '/ai',
    name: 'Ai',
    component: Ai
  },
  {
    path: '/dashboards',
    name: 'Dashboards',
    component: Dashboards
  },
  {
    path: '/dashboards/edit',
    name: 'DashboardEdit',
    component: DashboardEdit
  },
  {
    path: '/dashboards/view',
    name: 'DashboardView',
    component: DashboardView
  },
  {
    path: '/device-profiles',
    name: 'DeviceProfiles',
    component: DeviceProfiles
  },
  {
    path: '/ingest-logs',
    name: 'IngestLogs',
    component: IngestLogs
  },
  {
    path: '/data-browser',
    name: 'DataBrowser',
    component: DataBrowser
  },
  {
    path: '/settings',
    name: 'Settings',
    component: Settings
  },
  {
    path: '/extensions',
    name: 'Extensions',
    component: ExtensionsHub
  },
  {
    path: '/extensions/:name',
    name: 'ExtensionDetail',
    component: ExtensionDetail
  },
  {
    path: '/login',
    name: 'Login',
    component: Login
  }
];

const router = createRouter({
  history: createWebHistory(),
  routes
});

router.beforeEach(async (to, from, next) => {
  const { isAuthenticated, waitForInit } = useAuth()

  await waitForInit()

  if (to.name === 'Login') {
    if (isAuthenticated.value) {
      next('/')
    } else {
      next()
    }
  } else {
    if (!isAuthenticated.value) {
      next({ name: 'Login', query: { redirect: to.fullPath } })
    } else {
      next()
    }
  }
})

export default router;
