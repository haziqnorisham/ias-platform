import { createRouter, createWebHistory } from 'vue-router';

// route components
import Home from '../views/Home.vue';
import About from '../views/About.vue';
import Devices from '../views/Devices.vue';
import Diagnostics from '../views/Diagnostics.vue';
import Ai from '../views/Ai.vue';
import Dashboards from '../views/Dashboards.vue';
import DeviceProfiles from '../views/DeviceProfiles.vue';
import IngestLogs from '../views/IngestLogs.vue';

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
    path: '/device-profiles',
    name: 'DeviceProfiles',
    component: DeviceProfiles
  },
  {
    path: '/ingest-logs',
    name: 'IngestLogs',
    component: IngestLogs
  }
];

const router = createRouter({
  history: createWebHistory(),
  routes
});

export default router;
