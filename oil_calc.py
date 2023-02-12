from calculation import Calculation

class OilCalculation(Calculation):

    def calculate(self):
        result = []
        try:
            if float(self.quantity) > 1:
                emissions = round(42.8 * 0.074 * 0.845 / 1000 * float(self.quantity.replace(",", ".")), 2)
            else:
                emissions = round(42.8 * 0.074 * 0.845 / 1000 * float(self.quantity.replace(",", ".")), 3)
            emission_component_result = round(emissions, 2) * 30 * 1.19
            energy_content = 42.8 * 0.845 / 1000 * float(self.quantity.replace(",", "."))
            result.append(True)
            result.append(emissions)
            result.append(round(emission_component_result, 2))
            result.append(round(energy_content, 3))
        except:
            result.append(False)
        return result
        