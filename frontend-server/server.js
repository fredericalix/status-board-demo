// Charger les variables d'environnement
require("dotenv").config();

// Importer les modules nécessaires
const express = require("express");
const history = require("connect-history-api-fallback");

// Créer une nouvelle application Express
const app = express();

// Synchroniser le mode history de Vue.js avec l'application Express
app.use(
  history({
    // Activer les logs
    verbose: true,
  })
);

// Servir tous les fichiers du répertoire '/dist'
app.use(express.static("dist"));

// Fonction middleware pour forcer la redirection vers HTTPS
const forceHttps = (req, res, next) => {
  if (req.secure || req.headers["x-forwarded-proto"] === "https") {
    return next();
  } else {
    return res.redirect("https://" + req.headers.host + req.url);
  }
};

// Utiliser la fonction middleware pour forcer la redirection HTTPS
app.use(forceHttps);

// Lancer l'application sur le port spécifié par la variable d'environnement
app.listen(process.env.PORT, "0.0.0.0", () => {
  console.log("Lancement de l'application Vue.js et elle écoute sur le port " + process.env.PORT);
});
