import re

class Calculation:

    def __init__(self, quantity : str):
        self.quantity = quantity
    
    def format_output(self, input : float, is_energy_content : bool = False) -> str:
        if is_energy_content:
            return re.sub(r'(?<!^)(?=(\d{3})+,)', r'.',str(round(input, 3)).replace(".", ","))
        else:
            return re.sub(r'(?<!^)(?=(\d{3})+,)', r'.',str(format(input, '.2f')).replace(".", ","))