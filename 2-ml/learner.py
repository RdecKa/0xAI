class Learner:

    def __init__(self, models, feature_names):
        self.models = models
        self.feature_names = feature_names

    def get_models(self):
        return self.models
