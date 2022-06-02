<template>
  <h2>Connected Users</h2>
  <ul class="list-group" v-for="user in users" v-bind:key="user.id">
    <li class="list-group-item">
      Client <strong>#{{ user.id }}</strong>
      <span class="at">@</span>
      <span class="channel">{{ user.channel }}</span>
    </li>
  </ul>
</template>

<script>
export default {
  name: 'AdminUsers',
  props: {
    users: []
  }
};
</script>

<style scoped>
li > span.at {
  font-weight: bold;
  color: #009688;
}

li > span.channel {
  font-weight: bold;
  color: #795548;
}
</style>