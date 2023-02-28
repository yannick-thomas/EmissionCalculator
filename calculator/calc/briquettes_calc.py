from calculator.calc.calculation import Calculation

class BriquettesCalculation(Calculation):

    def calculate(self):
        result = []
        try:
            emissions = 19 * 0.0992 * 1 * 1000 * float(
                self.quantity.replace(",", ".")
            )
            emission_component_result =  emissions * 30 * 1.19 / 1000
            energy_content = 19 * 1 * float(self.quantity.replace(",", "."))
            result.append(True)
            result.append(emissions) 
            result.append(emission_component_result) 
            result.append(energy_content)
        except:
            result.append(False)
        return result
        