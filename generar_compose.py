import sys
import yaml

class IndentedDumper(yaml.Dumper):
    """Custom YAML dumper to correctly indent lists."""
    def increase_indent(self, flow=False, indentless=False):
        return super(IndentedDumper, self).increase_indent(flow, False)

def generar_compose(archivo_salida, cantidad_clientes):
    docker_compose = {
        'name': 'tp0',
        'services': {
            'server': {
                'container_name': 'server',
                'image': 'server:latest',
                'volumes': [
                    './server/config.ini:/config.ini'
                ],
                'entrypoint': 'python3 /main.py',
                'environment': [
                    'PYTHONUNBUFFERED=1'
                ],
                'networks': [
                    'testing_net'
                ]
            },
            'unzipper': {
                'image': 'busybox',
                'volumes': [
                    './.data:/data'
                ],
                'entrypoint': 'unzip -n /data/dataset.zip -d /data',
            }
        },
        'networks': {
            'testing_net': {
                'ipam': {
                    'driver': 'default',
                    'config': [
                        {'subnet': '172.25.125.0/24'}
                    ]
                }
            }
        }
    }

    for i in range(1, cantidad_clientes + 1):
        docker_compose['services'][f'client{i}'] = {
            'container_name': f'client{i}',
            'image': 'client:latest',
            'volumes': [
                './client/config.yaml:/config.yaml',
                f'./.data/agency-{i}.csv:/data/agency-{i}.csv:ro'
            ],
            'entrypoint': '/client',
            'environment': [
                f'CLI_ID={i}',
            ],
            'networks': [
                'testing_net'
            ],
            'depends_on': [
                'server',
                'unzipper'
            ]
        }

    with open(archivo_salida, 'w') as archivo:
        yaml.dump(docker_compose, archivo, Dumper=IndentedDumper, default_flow_style=False, sort_keys=False, indent=2)

if __name__ == "__main__":
    if len(sys.argv) != 3:
        sys.exit(1)

    archivo_salida = sys.argv[1]
    cantidad_clientes = int(sys.argv[2])
    generar_compose(archivo_salida, cantidad_clientes)
