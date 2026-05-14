<script setup>
// App.vue is now just a layout wrapper; content is rendered by router-view.
import SideNav from "./components/SideNav.vue";
import Breadcrumb from 'primevue/breadcrumb';
import ScrollPanel from 'primevue/scrollpanel';
import Footer from "./components/Footer.vue";
import Toast from 'primevue/toast';
import { ref, watch } from "vue";
import { useRoute } from 'vue-router';

const route = useRoute();

const home = ref({
  icon: 'pi pi-home'
});

// Map route names to breadcrumb labels
const routeLabelMap = {
  Home: 'Home',
  Ai: 'IAS AI (Preview)',
  Dashboards: 'Dashboards',
  DashboardEdit: 'Edit Dashboard',
  Devices: 'Devices',
  DeviceProfiles: 'Device Profiles',
  IngestLogs: 'Ingest Logs',
  Diagnostics: 'Diagnostics',
  Settings: 'Settings',
  About: 'About',
};

const items = ref([
  { label: 'Home', to: '/' }
]);

// Update breadcrumbs whenever the route changes.
// The :home prop already renders the home icon, so items should only
// contain the subsequent path segments (no "Home" entry).
watch(
  () => route.path,
  (newPath) => {
    const pathSegments = newPath.split('/').filter(Boolean);
    
    if (pathSegments.length === 0) {
      items.value = [];
      return;
    }

    const crumbs = [];
    let accumulatedPath = '';
    
    for (const segment of pathSegments) {
      accumulatedPath += '/' + segment;
      const routeName = route.matched.find(r => r.path === accumulatedPath)?.name;
      const label = routeLabelMap[routeName] || segment.charAt(0).toUpperCase() + segment.slice(1).replace(/-/g, ' ');
      
      crumbs.push({
        label: label,
        to: accumulatedPath
      });
    }

    items.value = crumbs;
  },
  { immediate: true }
);

</script>

<template>
  <Toast position="top-center"/>
  <SideNav />

  <main class="main">
    <Breadcrumb :home="home" :model="items" class="breadcrumb" />
    <div class="routerView">
      <router-view />
    </div>
    
    <ScrollPanel class="content">

    </ScrollPanel>
    
    <div class="footer_position">
      <Footer />
    </div>
  </main>

</template>

<style scoped>
.main {
  position: fixed;
  top: 0;
  right: 0;
  bottom: 0;
  left: 16rem;
  display: flex;
  flex-direction: column;
  padding-left: 40px;
  padding-right: 40px;
  box-sizing: border-box;
  overflow-y: auto;
  overflow-x: hidden;
}

.content {
  flex: 1;
  min-height: 0;
  border-radius: 12px;
  border: 0px solid #212121;
}
</style>

<style scoped>
.cards_container {
  margin: 10px;
  display: flex;
  flex-direction: row;
  flex-wrap: wrap;
  justify-content: space-evenly;
  flex-wrap: wrap;
}

.main_container {
  margin-left: 20px;
  margin-right: 20px;
  display: flex;
  flex-direction: row;
  flex-wrap: wrap;
  justify-content: space-evenly;
}

.breadcrumb {
  border-radius: 12px;
  margin-bottom: 20px;
  background-color: transparent;
}
</style>
