from calculation import Calculation

class OilCalculation(Calculation):

    def calculate(self):
        result = []
        try:
            emissions = 42.8 * 0.074 * 0.845 * float(self.quantity.replace(",", "."))
            emission_component_result = emissions * 30 * 1.19
            result.append(emissions)
            result.append(emission_component_result)
        except:
            return
        return result
        