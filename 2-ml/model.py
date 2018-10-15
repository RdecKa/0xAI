class Model:

    def __init__(self, model, feature_names):
        self.model = model
        self.feature_names = feature_names
        self.ID = ""

    def get_ID(self):
        return self.ID

    def fit(self, X, y):
        self.model.fit(X, y)

    def predict(self, X):
        return self.model.predict(X)

    def score(self, X, y):
        return self.model.score(X, y)
