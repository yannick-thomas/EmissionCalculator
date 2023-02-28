from calculator.calc.calculation import Calculation

class OilCalculation(Calculation):

    def calculate(self):
        result = []
        try:
            emissions = 42.8 * 0.074 * 0.845 * float(self.quantity.replace(",", "."))
            emission_component_result = round(emissions, 2) / 1000 * 30 * 1.19
            energy_content = 42.8 * 0.845 / 1000 * float(self.quantity.replace(",", ".")) * 277.778
            result.append(True)
            result.append(self.format_output(emissions))
            result.append(self.format_output(emission_component_result))
            result.append(self.format_output(energy_content, True
        except:
            result.append(False)
        return result
        