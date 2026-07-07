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
import IntegrationsHub from '../views/IntegrationsHub.vue';
import OnvifStreams from '../views/OnvifStreams.vue';
import OnvifStreamDetail from '../views/OnvifStreamDetail.vue';
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
    path: '/integrations',
    name: 'Integrations',
    component: IntegrationsHub
  },
  {
    path: '/integrations/onvif-streams',
    name: 'OnvifStreams',
    component: OnvifStreams
  },
  {
    path: '/integrations/onvif-streams/:id',
    name: 'OnvifStreamDetail',
    component: OnvifStreamDetail
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

router.beforeEach((to, from, next) => {
  const { isAuthenticated } = useAuth()

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
