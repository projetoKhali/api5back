import pandas as pd

# Ler todas as sheets da planilha Excel em DataFrames
sheet = './etl/dados_dw.xlsx'

dfs = pd.read_excel(sheet, sheet_name=None)  # Carrega todas as abas como dicionário de DataFrames

# Separar cada DataFrame por nome de aba
df_department = dfs['Departamentos']
df_user = dfs['Usuários']
df_process = dfs['Processos']
df_vacancy = dfs['Vagas']
df_candidates = dfs['Candidatos']
df_vacancy_candidate = dfs['Vaga-Candidatos']
df_interview = dfs['Entrevistas']
df_feedback = dfs['Feedbacks']
df_hiring = dfs['Contratações']

# print("Dados extraídos com sucesso!")
# print(df_vacancy.head())  # Exibe as primeiras linhas de um DataFrame para confirmar
