<template>
  <AppToolbar />
  <AdminHome />
</template>

<script>
import AppToolbar from '@/components/Toolbar';
import AdminHome from '@/components/Home';

export default {
  name: 'App',
  components: {
    AppToolbar,
    AdminHome
  }
};
</script>
