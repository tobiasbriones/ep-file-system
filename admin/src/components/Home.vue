<template>
  <div class="px-5 py-4">
    <AdminUsers :users="users" />
  </div>
</template>

<script>
import AdminUsers from '@/components/Users';

export default {
  name: 'AdminHome',
  components: { AdminUsers },
  data() {
    return {
      users: [
        {
          id: 1,
          channel: 'test'
        }
      ]
    };
  }
};
</script>
