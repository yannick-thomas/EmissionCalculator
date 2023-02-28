from calculator.calc.calculation import Calculation

class BriquettesCalculation(Calculation):
    
    def calculate(self):
        result = []
        try:
            emissions = 19 * 0.0992 * 1 * float(
                self.quantity.replace(",", ".")
            )
            emission_component_result =  emissions * 30 * 1.19
            energy_content = 19 * 1 * float(self.quantity.replace(",", "."))

            result.append(True)
            result.append(self.format_output(emissions)) 
            result.append(self.format_output(emission_component_result)) 
            result.append(self.format_output(energy_content, True))

        except:
            result.append(False)

        return result
        