<template>
  <div id="app">
    <div class="grid-container">
      <div
        v-for="status in statusList"
        :key="status.id"
        class="status-box"
        :style="{ backgroundColor: status.state }"
      >
        <p>{{ status.designation }}</p>
      </div>
    </div>
  </div>
</template>

<script>
import axios from "axios";

const API_URL = "http://localhost:8080"; // Remplacez par l'URL de votre backend

export default {
  data() {
    return {
      statusList: [],
    };
  },
  async created() {
    try {
      const response = await axios.get(`${API_URL}/status`);
      this.statusList = response.data;
      this.setupEventSource();
    } catch (error) {
      console.error("Erreur lors de la récupération des données :", error);
    }
  },
  methods: {
    setupEventSource() {
  const eventSource = new EventSource(`${API_URL}/events`);
  eventSource.onmessage = () => { // Supprimez 'event' ici
    this.fetchStatus();
  };
},
    async fetchStatus() {
      try {
        const response = await axios.get(`${API_URL}/status`);
        this.statusList = response.data;
      } catch (error) {
        console.error("Erreur lors de la récupération des données :", error);
      }
    },
  },
};
</script>

<style>
body {
  background-color: #1e1e1e;
  margin: 0;
  font-family: "Segoe UI", Tahoma, Geneva, Verdana, sans-serif;
}

.grid-container {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
  padding: 20px;
}

.status-box {
  background-color: #eee;
  border-radius: 10px;
  padding: 20px;
  text-align: center;
  color: white;
}
</style>
