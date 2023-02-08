from calculation import Calculation

class BriquettesCalculation(Calculation):

    def calculate(self):
        result = []
        try:
            emissions = 19 * 0.0992 * 1 * float(
                self.quantity.replace(",", ".")
            )
            emission_component_result =  emissions * 30 * 1.19
            energy_content = 19 * 1 * float(self.quantity.replace(",", "."))
            result.append(emissions, emission_component_result, energy_content)
        except:
            result.append('Fehler')
        return result
        