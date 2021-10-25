<template>
  <div>
    <el-form ref="form" label-width="120px">
      <el-form-item label="Место">
        <el-select class="el-select"
                   v-model="value"
                   filterable
                   remote
                   reserve-keyword
                   placeholder="Регион, город, улица (напр. Цветной проезд)"
                   :remote-method="remoteMethod"
                   :loading="loading"
                   :loading-text="loading_text">
          <el-option
              v-for="item in options"
              :key="item.point.lat+item.point.lng"
              :label="item.country +', '+item.state+', '+item.name"
              :value="item.point.lat+item.point.lng">
            <span style="float: left">{{ item.country + ', ' + item.state + ', ' + item.name }}</span>
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="Радиус, м">
        <el-input-number v-model="radius" controls-position="right" :min="0" :step="100"></el-input-number>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="onSubmit">Поиск</el-button>
      </el-form-item>
    </el-form>

    <el-descriptions title="Погода" v-if="weather">
      <el-descriptions-item label="Состояние">
        <el-tag size="small">{{ weather.weather[0].description }}</el-tag>
      </el-descriptions-item>
      <el-descriptions-item label="Температура">{{ (parseFloat(weather.main.temp) - 273.1).toFixed(1)}}℃</el-descriptions-item>
      <el-descriptions-item label="Ощущается как">{{ (parseFloat(weather.main.feels_like) -273.1).toFixed(1)}}℃</el-descriptions-item>
      <el-descriptions-item label="Влажность">{{ weather.main.humidity }}%</el-descriptions-item>
      <el-descriptions-item label="Давление">{{ weather.main.pressure }}гПа</el-descriptions-item>
      <el-descriptions-item label="Скорость ветра">{{ weather.wind.speed }}м/с</el-descriptions-item>
      <el-descriptions-item label="Направление ветра">{{ weather.wind.deg }}°</el-descriptions-item>
    </el-descriptions>

    <el-descriptions title="Места" v-if="places">
    </el-descriptions>

    <el-table
        :data="places"
        style="width: 100%"
        v-if="places">
      <el-table-column
          prop="properties.name"
          label="Название">
      </el-table-column>
      <el-table-column
          fixed="right"
          label="Детали"
          width="150">
        <template slot-scope="scope">
          <el-button @click="showDetail(scope.$index)">Детали</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog
        title="Детали"
        :visible.sync="detailsDialogVisible"
        v-if="detailsDialogVisible"
        width="80%">
      <img :src="description.image" alt="No image)" class="image"/>
      <div v-html="description_info"></div>
      <div>
        <a :href="description.otm">Open street maps</a>
      </div>
      <div>
        <a :href="description.wikipedia">Wikipedia</a>
      </div>
      <span slot="footer" class="dialog-footer">
        <el-button type="primary" @click="detailsDialogVisible = false">ОК</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  name: "Place",
  data() {
    return {
      value: null,
      options: [],
      error: "",
      loading: false,
      loading_text: "Loading",
      radius: 1000,
      places: null,
      description: null,
      weather: null,

      detailsDialogVisible: false,
    }
  },
  methods: {
    async remoteMethod(query) {
      this.loading = true;
      if (query !== '') {
        try {
          const response = await axios.get("/geocode-async/" + query)

          this.options = response.data.hits
          console.log(this.options)
        } catch (e) {
          this.error = "Error occurred!"
        }
      } else {
        this.options = []
      }
      this.loading = false;
    },
    async getPlaces(item) {
      const response_places = await axios.get("/places-async", {
        params: {
          radius: this.radius,
          lng: item.point.lng,
          lat: item.point.lat,
        }
      })

      this.places = response_places.data.features.filter(function (e) {
        return e.properties.name !== ""
      })
      console.log("Places")
      console.log(this.places)
    },
    async getWeather(item) {
      const response_weather = await axios.get("/weather-async", {
        params: {
          lng: item.point.lng,
          lat: item.point.lat,
        }
      })

      this.weather = response_weather.data
      console.log("Weather")
      console.log(this.weather)
    },
    async onSubmit() {
      for (const i in this.options) {
        let item = this.options[i]
        if (item.point.lng + item.point.lat === this.value) {
          console.log(item.point)

          try {
            this.getPlaces(item)
            this.getWeather(item)
          } catch (e) {
            console.log("Error")
            console.log(e)
            this.error = "Error occurred!"
          }

          break;
        }
      }
    },
    async showDetail(idx) {
      let xid = this.places[idx].properties.xid;
      try {
        const response = await axios.get("/description/" + xid)

        this.description = response.data

        console.log("Description")
        console.log(this.description)

        this.detailsDialogVisible = true
      } catch (e) {
        console.log("Error")
        console.log(e)
        this.error = "Error occurred!"
      }
    }
  },
  computed: {
    // eslint-disable-next-line vue/return-in-computed-property
    description_info: function () {
      if (this.description != null) {
        if (this.description.info != null) {
          return this.description.info.descr
        }
      } else {
        return "Nothing"
      }
    }
  }
}
</script>

<style scoped>

.el-select {
  width: 450px;
}

.image {
  width: 600px;
}

</style>