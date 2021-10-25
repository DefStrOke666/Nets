<template>
  <div>
    <el-tabs class="app"  type="border-card">
      <el-tab-pane label="Info">
        <Place class="place"/>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script>
import Place from "@/components/Place";
export default {
  name: 'App',
  components: {
    Place
  }
}
</script>

<style>
.app {
  width: 700px;
  margin: 100px auto 0;
}

</style>
