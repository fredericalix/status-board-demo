<template>
  <div id="app">
    <!-- Création d'un conteneur pour les notifications -->
    <div class="notification-container">
      <!-- Boucle sur les notifications pour créer un carré pour chaque notification -->
      <div
        v-for="(notification, index) in notifications"
        :key="index"
        class="notification"
        :style="{ backgroundColor: notification.state }"
      >
        <!-- Affichage de la désignation de la notification au centre du carré -->
        <p>{{ notification.designation }}</p>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios';

// Définition de l'URL de l'API backend en fonction de la variable d'environnement ou par défaut en utilisant "http://localhost:8080"
const API_URL = process.env.API_URL || 'http://localhost:8080';

export default {
  name: 'App',
  data() {
    return {
      notifications: [],
    };
  },
  async created() {
    // Récupération des notifications depuis l'API lors du chargement de l'application frontend
    try {
      const response = await axios.get(`${API_URL}/status`);
      this.notifications = response.data;
    } catch (error) {
      console.error('Erreur lors de la récupération des notifications:', error);
    }

    // Écoute des mises à jour des notifications via Server-Sent Events (SSE)
    const eventSource = new EventSource(`${API_URL}/events`);
    eventSource.onmessage = (event) => {
      const updatedNotification = JSON.parse(event.data);
      this.updateNotification(updatedNotification);
    };
  },
  methods: {
    // Mise à jour d'une notification dans la liste des notifications
    updateNotification(updatedNotification) {
      this.notifications = this.notifications.map((notification) => {
        if (notification.id === updatedNotification.id) {
          return updatedNotification;
        }
        return notification;
      });
    },
  },
};
</script>

<style scoped>
body {
  background-color: #2c3e50;
  font-family: 'Source Sans Pro', sans-serif;
}

.notification-container {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  align-items: center;
  gap: 1rem;
  padding: 2rem;
}

.notification {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100px;
  height: 100px;
  border-radius: 8px;
  font-size: 1.2rem;
  color: #fff;
  text-align: center;
}
</style>
