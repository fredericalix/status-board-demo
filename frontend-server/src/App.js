// Importer les fonctions nécessaires depuis les packages 'react' et 'axios'
import React, { useState, useEffect } from "react";
import axios from "axios";

// Importer le fichier de style CSS
import "./index.css";

// Définir l'URL de base de l'API en utilisant la variable d'environnement 'REACT_APP_API_URL'
const API_BASE_URL = process.env.REACT_APP_API_URL;

// Définition du composant 'App'
const App = () => {
  // Initialiser 'statusList' à un tableau vide et 'setStatusList' à une fonction permettant de le mettre à jour
  const [statusList, setStatusList] = useState([]);

  // Définir une fonction asynchrone 'fetchData' pour récupérer des données depuis l'API
  const fetchData = async () => {
    try {
      // Faire une requête GET vers l'endpoint '/status' de l'API et attendre la réponse
      const response = await axios.get(`${API_BASE_URL}/status`);

      // Mettre à jour 'statusList' avec les données reçues dans la réponse
      setStatusList(response.data);
    } catch (error) {
      // En cas d'erreur, afficher un message d'erreur dans la console
      console.error("Erreur lors de la récupération des données :", error);
    }
  };

  // Définir une fonction 'setupEventSource' pour configurer une connexion EventSource
  const setupEventSource = () => {
    // Créer une nouvelle connexion EventSource vers l'endpoint '/events' de l'API
    const eventSource = new EventSource(`${API_BASE_URL}/events`);

    // Lorsqu'un message est reçu via la connexion EventSource, appeler 'fetchData'
    eventSource.onmessage = (e) => {
      fetchData();
    };

    // En cas d'erreur avec la connexion EventSource, afficher un message d'erreur dans la console
    eventSource.onerror = (e) => {
      console.error("Erreur lors de la connexion à l'EventSource:", e);
    };
  };

  // Lorsque le composant est monté, appeler 'fetchData' et 'setupEventSource'
  useEffect(() => {
    fetchData();
    setupEventSource();
  }, []);

  // Rendre le composant
  return (
    <div className="container">
      {/* Parcourir 'statusList' et pour chaque 'status', créer une div avec l'id en tant que clé, le 'state' en tant que couleur de fond, et la 'designation' comme contenu */}
      {statusList.map((status) => (
        <div
          key={status.id}
          className="status-box"
          style={{ backgroundColor: status.state }}
        >
          {status.designation}
        </div>
      ))}
    </div>
  );
};

// Exporter le composant 'App' pour pouvoir l'utiliser dans d'autres fichiers
export default App;